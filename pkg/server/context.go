package server

import (
	"context"

	"github.com/k3s-io/k3s/pkg/util"
	"github.com/k3s-io/k3s/pkg/version"
	"github.com/rancher/wrangler/v3/pkg/generated/controllers/apps"
	"github.com/rancher/wrangler/v3/pkg/generated/controllers/batch"
	"github.com/rancher/wrangler/v3/pkg/generated/controllers/core"
	"github.com/rancher/wrangler/v3/pkg/generated/controllers/rbac"
	"github.com/rancher/wrangler/v3/pkg/start"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
)

type Context struct {
	Batch  *batch.Factory
	Apps   *apps.Factory
	Auth   *rbac.Factory
	Core   *core.Factory
	K8s    kubernetes.Interface
	Event  record.EventRecorder
	Config *rest.Config
}

func (c *Context) Start(ctx context.Context) error {
	return start.All(ctx, 5, c.Apps, c.Auth, c.Batch, c.Core)
}

func NewContext(ctx context.Context, config *Config, forServer bool) (*Context, error) {
	cfg := config.ControlConfig.Runtime.KubeConfigAdmin
	if forServer {
		cfg = config.ControlConfig.Runtime.KubeConfigSupervisor
	}
	restConfig, err := clientcmd.BuildConfigFromFlags("", cfg)
	if err != nil {
		return nil, err
	}
	restConfig.UserAgent = util.GetUserAgent(version.Program + "-supervisor")

	k8s, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	var recorder record.EventRecorder
	if forServer {
		recorder = util.BuildControllerEventRecorder(k8s, version.Program+"-supervisor", metav1.NamespaceAll)
	}

	return &Context{
		K8s:    k8s,
		Auth:   rbac.NewFactoryFromConfigOrDie(restConfig),
		Apps:   apps.NewFactoryFromConfigOrDie(restConfig),
		Batch:  batch.NewFactoryFromConfigOrDie(restConfig),
		Core:   core.NewFactoryFromConfigOrDie(restConfig),
		Event:  recorder,
		Config: restConfig,
	}, nil
}
