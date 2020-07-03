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

func TestGKE(t *testing.T) {
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
	err = ioutil.WriteFile(path.Join(clusterDir, "cluster.tf"), []byte(gkeClusterTF), 0644)
	if err != nil {
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
			"signing_cert": signingCert,
			"signing_key":  signingKey,
		},
	}

	// Run `terraform init` and `terraform apply`
	result = terraform.InitAndApply(t, platformTFOptions)
	t.Logf("Platform setup stdout: \n%v", result)
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

const gkeClusterTF = `
module "cluster" {
  source   = "github.com/telliott-io/kube-clusters//gke?ref=v0.3.0"
  cluster_name = var.cluster_name
}

provider "google" {
  credentials = base64decode(var.gcloud_credentials_base64)
  project     = "telliott-io"
  region      = "us-central1"
  zone        = "us-central1-c"
}

variable "gcloud_credentials_base64" {}

variable "cluster_name" {}

output "config" {
	value = module.cluster.kubernetes
	sensitive = true
}
`
