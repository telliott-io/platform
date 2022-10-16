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

func TestDigitalOcean(t *testing.T) {
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
	err = ioutil.WriteFile(path.Join(clusterDir, "cluster.tf"), []byte(clusterTF), 0644)
	if err != nil {
		t.Fatal(err)
	}

	clusterTFOptions := &terraform.Options{
		TerraformDir: clusterDir,
		// EnvVars: map[string]string{
		// 	"TF_LOG": "debug",
		// },
	}

	// At the end of the test, run `terraform destroy`
	defer terraform.Destroy(t, clusterTFOptions)

	// Run `terraform init` and `terraform apply`
	result := terraform.InitAndApply(t, clusterTFOptions)
	config := terraform.OutputMap(t, clusterTFOptions, "config")
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

	platformTFOptions := &terraform.Options{
		TerraformDir: platformDir,
		Vars: map[string]interface{}{
			"kubernetes": string(kubernetesJSON),
		},
	}

	// Run `terraform init` and `terraform apply`
	result = terraform.InitAndApply(t, platformTFOptions)
	t.Logf("Platform setup stdout: \n%v", result)
}

const clusterTF = `
terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      version = "2.5.1"
    }
  }
}

module "cluster" {
  source   = "github.com/telliott-io/kube-clusters//digitalocean?ref=v0.6.1"
  cluster_name = var.cluster_name
}
provider "digitalocean" {
  token = var.do_token
}
variable "do_token" {}
variable "cluster_name" {}
output "config" {
	value = module.cluster.kubernetes
	sensitive = true
}
`

const platformTF = `
module "platform" {
	source   = "../../"
	kubernetes = var.kubernetes
	environment = "platform-test"
	hostname = "platform.test"
	argocd_admin_password = "secret"
	bootstrap_repository = "https://telliott-io.github.io/testbootstrap"
	bootstrap_chart = "bootstrap"
	bootstrap_version = "0.1.1"
}

variable kubernetes {}
`
