package cloudprovider

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/k3s-io/k3s/pkg/version"
	cloudprovider "k8s.io/cloud-provider"
)

// Config describes externally-configurable cloud provider configuration.
// This is normally unmarshalled from a JSON config file.
type Config struct {
	NodeEnabled bool `json:"nodeEnabled"`
	Rootless    bool `json:"rootless"`
}

type k3s struct {
	Config
}

var _ cloudprovider.Interface = &k3s{}

func init() {
	cloudprovider.RegisterCloudProvider(version.Program, func(config io.Reader) (cloudprovider.Interface, error) {
		var err error
		k := k3s{
			Config: Config{
				NodeEnabled: true,
			},
		}

		if config != nil {
			var bytes []byte
			bytes, err = io.ReadAll(config)
			if err == nil {
				err = json.Unmarshal(bytes, &k.Config)
			}
		}

		if !k.NodeEnabled {
			return nil, fmt.Errorf("all cloud-provider functionality disabled by config")
		}

		return &k, err
	})
}

func (k *k3s) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {}

func (k *k3s) Instances() (cloudprovider.Instances, bool) {
	return nil, false
}

func (k *k3s) InstancesV2() (cloudprovider.InstancesV2, bool) {
	return k, k.NodeEnabled
}

func (k *k3s) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return nil, false
}

func (k *k3s) Zones() (cloudprovider.Zones, bool) {
	return nil, false
}

func (k *k3s) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

func (k *k3s) Routes() (cloudprovider.Routes, bool) {
	return nil, false
}

func (k *k3s) ProviderName() string {
	return version.Program
}

func (k *k3s) HasClusterID() bool {
	return false
}
