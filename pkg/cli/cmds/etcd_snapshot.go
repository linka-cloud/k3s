package cmds

import (
	"github.com/k3s-io/k3s/pkg/version"
	"github.com/urfave/cli"
)

const EtcdSnapshotCommand = "etcd-snapshot"

var EtcdSnapshotFlags = []cli.Flag{
	DebugFlag,
	ConfigFlag,
	LogFile,
	AlsoLogToStderr,
	&cli.StringFlag{
		Name:        "node-name",
		Usage:       "(agent/node) Node name",
		EnvVar:      version.ProgramUpper + "_NODE_NAME",
		Destination: &AgentConfig.NodeName,
	},
	DataDirFlag,
	&cli.StringFlag{
		Name:        "etcd-token,t",
		Usage:       "(cluster) Shared secret used to authenticate to etcd server",
		Destination: &ServerConfig.Token,
	},
	&cli.StringFlag{
		Name:        "etcd-server, s",
		Usage:       "(cluster) Server with etcd role to connect to for snapshot management operations",
		Value:       "https://127.0.0.1:6443",
		Destination: &ServerConfig.ServerURL,
	},
}
