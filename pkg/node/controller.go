package node

import (
	"context"

	"github.com/pkg/errors"
	coreclient "github.com/rancher/wrangler/v3/pkg/generated/controllers/core/v1"
	"github.com/sirupsen/logrus"
	core "k8s.io/api/core/v1"

	"github.com/k3s-io/k3s/pkg/nodepassword"
)

func Register(ctx context.Context,
	secrets coreclient.SecretController,
	configMaps coreclient.ConfigMapController,
	nodes coreclient.NodeController,
) error {
	h := &handler{
		secrets:    secrets,
		configMaps: configMaps,
	}
	nodes.OnChange(ctx, "node", h.onChange)
	nodes.OnRemove(ctx, "node", h.onRemove)

	return nil
}

type handler struct {
	secrets    coreclient.SecretController
	configMaps coreclient.ConfigMapController
}

func (h *handler) onChange(key string, node *core.Node) (*core.Node, error) {
	if node == nil {
		return nil, nil
	}
	return h.updateHosts(node, false)
}

func (h *handler) onRemove(key string, node *core.Node) (*core.Node, error) {
	return h.updateHosts(node, true)
}

func (h *handler) updateHosts(node *core.Node, removed bool) (*core.Node, error) {
	var nodeName string
	nodeName = node.Name
	if removed {
		if err := h.removeNodePassword(nodeName); err != nil {
			logrus.Warn(errors.Wrap(err, "Unable to remove node password"))
		}
	}
	return nil, nil
}

func (h *handler) removeNodePassword(nodeName string) error {
	return nodepassword.Delete(h.secrets, nodeName)
}
