package kind

import (
	"fmt"
	"os"
	"strings"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cmd"
)

// New creates a KIND cluster for testing
func New(name, kubeConfigPath string) (cleanup func() error, err error) {
	kindProvider := cluster.NewProvider(
		cluster.ProviderWithLogger(
			cmd.NewLogger(),
		),
	)
	err = kindProvider.Create(
		name,
		cluster.CreateWithNodeImage("kindest/node:v1.22.13"),
		cluster.CreateWithKubeconfigPath(kubeConfigPath),
		cluster.CreateWithRawConfig([]byte(kindConfigFile)),
	)
	if err != nil {
		return nil, err
	}

	return func() error {
		var errstrings []string
		err1 := kindProvider.Delete(name, "")
		if err1 != nil {
			errstrings = append(errstrings, err1.Error())
		}
		err2 := os.Remove(kubeConfigPath)
		if err2 != nil {
			errstrings = append(errstrings, err2.Error())
		}
		return fmt.Errorf(strings.Join(errstrings, "\n"))
	}, nil
}

const kindConfigFile = `
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 32080
    hostPort: 32080
    listenAddress: "127.0.0.1"
    protocol: TCP
- role: worker
- role: worker
`
