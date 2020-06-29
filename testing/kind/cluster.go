package kind

import (
	"fmt"
	"os"
	"strings"
	"time"

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
		cluster.CreateWithNodeImage("kindest/node:v1.16.9"),
		cluster.CreateWithRetain(false),
		cluster.CreateWithWaitForReady(time.Duration(0)),
		cluster.CreateWithKubeconfigPath(kubeConfigPath),
		cluster.CreateWithDisplayUsage(true),
		cluster.CreateWithDisplaySalutation(true),
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
