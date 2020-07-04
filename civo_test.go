package platform

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/telliott-io/platform/testing/testdir"
)

func TestCivo(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	workingDir := "testdata"
	workingDirCleanup, err := testdir.New(workingDir)
	if err != nil {
		t.Fatal(err)
	}
	defer workingDirCleanup()

	clusterDir := path.Join(workingDir, "cluster")
	if err := os.Mkdir(clusterDir, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(path.Join(clusterDir, "cluster.tf"), []byte(civoClusterTF), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Copy provider
	if err := os.Mkdir(path.Join(clusterDir, "terraform.d"), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(path.Join(clusterDir, "terraform.d", "plugins"), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	providerDir := path.Join(clusterDir, "terraform.d", "plugins", "linux_amd64")
	if err := os.Mkdir(providerDir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := copyFile("terraform-provider-civo", path.Join(providerDir, "terraform-provider-civo")); err != nil {
		t.Fatal(err)
	}

	varFiles := []string{}
	if fileExists("terraform.tfvars") {
		varFiles = []string{"../../terraform.tfvars"}
	}

	clusterTFOptions := &terraform.Options{
		TerraformDir: clusterDir,
		VarFiles:     varFiles,
	}

	// At the end of the test, run `terraform destroy`
	defer terraform.Destroy(t, clusterTFOptions)

	// Run `terraform init` and `terraform apply`
	result := terraform.InitAndApply(t, clusterTFOptions)
	config := terraform.OutputMap(t, clusterTFOptions, "config")
	dnsName := terraform.Output(t, clusterTFOptions, "dns_name")
	t.Logf("Cluster creation stdout: \n%v", result)

	platformDir := path.Join(workingDir, "platform")
	if err := os.Mkdir(platformDir, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(path.Join(platformDir, "platform.tf"), []byte(platformTF), 0644)
	if err != nil {
		t.Fatal(err)
	}

	kubernetesJSON, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		t.Fatal(err)
	}

	signingCert, signingKey, err := createCerts()
	if err != nil {
		t.Fatal(err)
	}

	platformTFOptions := &terraform.Options{
		TerraformDir: platformDir,
		Vars: map[string]interface{}{
			"kubernetes":   string(kubernetesJSON),
			"dns_name":     dnsName,
			"signing_cert": signingCert,
			"signing_key":  signingKey,
		},
	}

	// Run `terraform init` and `terraform apply`
	result = terraform.InitAndApply(t, platformTFOptions)
	t.Logf("Platform setup stdout: \n%v", result)
}

func copyFile(src, dst string) error {
	content, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(dst, content, os.ModePerm); err != nil {
		return err
	}
	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.Chmod(dst, si.Mode()); err != nil {
		return err
	}
	return nil
}

const civoClusterTF = `
module "cluster" {
  source   = "github.com/telliott-io/kube-clusters//civo?ref=v0.4.0"
  cluster_name = var.cluster_name
}

provider "civo" {
  token = var.civo_api_key
}

variable "civo_api_key" {
}

variable "cluster_name" {}

output "config" {
	value = module.cluster.kubernetes
	sensitive = true
}

output "dns_name" {
	value = module.cluster.dns_name
}
`
