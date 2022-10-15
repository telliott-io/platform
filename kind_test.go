package platform

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/telliott-io/platform/testing/kind"
	"github.com/telliott-io/platform/testing/testdir"
)

func TestWithKind(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	tfdir := "testdata"
	workingDirCleanup, err := testdir.New(tfdir)
	if err != nil {
		t.Fatal(err)
	}
	defer workingDirCleanup()

	err = ioutil.WriteFile(path.Join(tfdir, "main.tf"), []byte(kindPlatformTF), 0644)
	if err != nil {
		t.Fatal(err)
	}

	kubeconfigfile := "kindconfig"
	kindCleanup, err := kind.New("argo-test-kind", path.Join(tfdir, kubeconfigfile))
	if err != nil {
		t.Fatal(err)
	}
	defer kindCleanup()

	terraformOptions := &terraform.Options{
		// The path to where your Terraform code is located
		TerraformDir: tfdir,
		Vars:         map[string]interface{}{},
		EnvVars: map[string]string{
			"TF_LOG": "debug",
		},
	}

	// At the end of the test, run `terraform destroy`
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`
	result := terraform.InitAndApply(t, terraformOptions)
	t.Logf("Stdout: %v", result)

}

const kindPlatformTF = `
module "platform" {
	source   = "../"
	kubernetes = "{\"config_path\": \"${path.module}/kindconfig\"}"
	environment = "platform-test"
	hostname = "platform.test"
	argocd_admin_password = "secret"
	bootstrap_repository = "https://telliott-io.github.io/testbootstrap"
	bootstrap_chart = "bootstrap"
	bootstrap_version = "0.1.1"

	debug = true
}
`
