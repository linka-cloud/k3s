//go:build !windows
// +build !windows

package rootlessports

import (
	"context"
	"time"

	"github.com/k3s-io/k3s/pkg/rootless"
	coreClients "github.com/rancher/wrangler/v3/pkg/generated/controllers/core/v1"
	"github.com/rootless-containers/rootlesskit/pkg/api/client"
	"github.com/rootless-containers/rootlesskit/pkg/port"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

var (
	all = "_all_"
)

func Register(ctx context.Context, serviceController coreClients.ServiceController, httpsPort int) error {
	var (
		err            error
		rootlessClient client.Client
	)

	if rootless.Sock == "" {
		return nil
	}

	for i := 0; i < 30; i++ {
		rootlessClient, err = client.New(rootless.Sock)
		if err == nil {
			break
		} else {
			logrus.Infof("Waiting for rootless API socket %s: %v", rootless.Sock, err)
			time.Sleep(1 * time.Second)
		}
	}
	if err != nil {
		return err
	}

	h := &handler{
		rootlessClient: rootlessClient,
		serviceClient:  serviceController,
		serviceCache:   serviceController.Cache(),
		httpsPort:      httpsPort,
		ctx:            ctx,
	}
	serviceController.OnChange(ctx, "rootlessports", h.serviceChanged)
	serviceController.Enqueue("", all)

	return nil
}

type handler struct {
	rootlessClient client.Client
	serviceClient  coreClients.ServiceController
	serviceCache   coreClients.ServiceCache
	httpsPort      int
	ctx            context.Context
}

func (h *handler) serviceChanged(key string, svc *v1.Service) (*v1.Service, error) {
	if key != all {
		h.serviceClient.Enqueue("", all)
		return svc, nil
	}

	ports, err := h.rootlessClient.PortManager().ListPorts(h.ctx)
	if err != nil {
		return svc, err
	}

	boundPorts := map[string]map[int]int{
		"tcp": {},
		"udp": {},
	}
	for _, port := range ports {
		boundPorts[port.Spec.Proto][port.Spec.ParentPort] = port.ID
	}

	toBindPort := map[string]map[int]int{
		"tcp": {h.httpsPort: h.httpsPort},
		"udp": {},
	}

	for proto, ports := range toBindPort {
		for bindPort, childBindPort := range ports {
			if _, ok := boundPorts[proto][bindPort]; ok {
				logrus.Debugf("Parent port %d/%s to child already bound", bindPort, proto)
				delete(boundPorts[proto], bindPort)
				continue
			}

			status, err := h.rootlessClient.PortManager().AddPort(h.ctx, port.Spec{
				Proto:      proto,
				ParentPort: bindPort,
				ChildPort:  childBindPort,
			})
			if err != nil {
				return svc, err
			}

			logrus.Infof("Bound parent port %s:%d/%s to child namespace port %d", status.Spec.ParentIP,
				status.Spec.ParentPort, proto, status.Spec.ChildPort)
		}
	}

	for proto, ports := range boundPorts {
		for bindPort, id := range ports {
			if err := h.rootlessClient.PortManager().RemovePort(h.ctx, id); err != nil {
				return svc, err
			}

			logrus.Infof("Removed parent port %d/%s to child namespace", bindPort, proto)
		}
	}

	return svc, nil
}
