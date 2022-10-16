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
	kindCleanup, err := kind.New("platform-test-kind", path.Join(tfdir, kubeconfigfile))
	if err != nil {
		t.Fatal(err)
	}
	defer kindCleanup()

	terraformOptions := &terraform.Options{
		TerraformDir: tfdir,
		Vars:         map[string]interface{}{},
		// EnvVars: map[string]string{
		// 	"TF_LOG": "debug",
		// },
	}

	// At the end of the test, run `terraform destroy`
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`
	_ = terraform.InitAndApply(t, terraformOptions)

	checkEnvironmentHost(t, "http://localhost:32080", "platform.test", "platform-test")
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
	service_type = "NodePort"
}
`
