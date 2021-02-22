package argocd

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"path"
	"strings"
	"testing"
	"time"

	sealedsecrets "github.com/bitnami-labs/sealed-secrets/pkg/apis/sealed-secrets/v1alpha1"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/telliott-io/platform/testing/kind"
	"github.com/telliott-io/platform/testing/testdir"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/deprecated/scheme"
	"k8s.io/client-go/util/cert"
)

func TestPlatform(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	tfdir := "testdata"
	workingDirCleanup, err := testdir.New(tfdir)
	if err != nil {
		t.Fatal(err)
	}
	defer workingDirCleanup()

	err = ioutil.WriteFile(path.Join(tfdir, "main.tf"), []byte(mainTF), 0644)
	if err != nil {
		t.Fatal(err)
	}

	kubeconfigfile := "kindconfig"
	kindCleanup, err := kind.New("platform-test-kind", path.Join(tfdir, kubeconfigfile))
	if err != nil {
		t.Fatal(err)
	}
	defer kindCleanup()

	signingCert, signingKey, err := createCerts()
	if err != nil {
		t.Fatal(err)
	}

	terraformOptions := &terraform.Options{
		// The path to where your Terraform code is located
		TerraformDir: tfdir,
		Vars: map[string]interface{}{
			"secret_signing_cert": signingCert,
			"secret_signing_key":  signingKey,
		},
	}

	// At the end of the test, run `terraform destroy`
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`
	result := terraform.InitAndApply(t, terraformOptions)
	t.Logf("Stdout: %v", result)

}

var secretObj = &v1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "mysecret",
		Namespace: "default",
	},
	Data: map[string][]byte{
		"foo": []byte(secretValue),
	},
}

func sealTestSecret(signingCert string) (*sealedsecrets.SealedSecret, error) {
	pubKey, err := parseKey(strings.NewReader(signingCert))
	if err != nil {
		return nil, err
	}
	return sealedsecrets.NewSealedSecret(scheme.Codecs, pubKey, secretObj)
}

func parseKey(r io.Reader) (*rsa.PublicKey, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	certs, err := cert.ParseCertsPEM(data)
	if err != nil {
		return nil, err
	}

	// ParseCertsPem returns error if len(certs) == 0, but best to be sure...
	if len(certs) == 0 {
		return nil, errors.New("Failed to read any certificates")
	}

	cert, ok := certs[0].PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("Expected RSA public key but found %v", certs[0].PublicKey)
	}

	return cert, nil
}

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

const secretValue = "bar"

const mainTF = `
module "platform" {
	source = "../"
	kubernetes = "{\"config_path\": \"${path.module}/kindconfig\"}"
	environment = "platform-test-1"
	hostname = "example.com"
	argocd_admin_password = "secret"
	secret_signing_cert = var.secret_signing_cert
    secret_signing_key = var.secret_signing_key
	bootstrap_repository = "https://telliott-io.github.io/bootstrap"
	bootstrap_chart = "bootstrap"
}

variable secret_signing_cert {}
variable secret_signing_key {}
`
