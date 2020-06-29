package platform

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"testing"
	"time"

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
		VarFiles:     []string{"../../terraform.tfvars"},
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

	t.Logf("Kubernetes json:\n%+v", string(kubernetesJSON))
	if true {
		return
	}

	signingCert, signingKey, err := createCerts()
	if err != nil {
		t.Fatal(err)
	}

	platformTFOptions := &terraform.Options{
		TerraformDir: platformDir,
		VarFiles:     []string{"../../terraform.tfvars"},
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

const clusterTF = `
module "cluster" {
  source   = "github.com/telliott-io/kube-clusters//digitalocean?ref=v0.3.0"
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

	secret_signing_cert = var.signing_cert
	secret_signing_key = var.signing_key

	environment = "platform-test"

	hostname = "platform.test"

	argocd_admin_password = "secret"

	bootstrap_repository = "https://telliott-io.github.io/bootstrap"
	bootstrap_chart = "bootstrap"
}

variable kubernetes {}
variable signing_cert {}
variable signing_key {}
`

func createCerts() (crt string, key string, err error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization: []string{"Some Org"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  false,
		ExtKeyUsage:           []x509.ExtKeyUsage{},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", err
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return "", "", err
	}
	caPEM := new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return "", "", err
	}

	caPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	if err != nil {
		return "", "", err
	}
	return string(caPEM.Bytes()), string(caPrivKeyPEM.Bytes()), nil
}
