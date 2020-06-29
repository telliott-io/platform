package secrets

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/telliott-io/platform/testing/kind"
	"github.com/telliott-io/platform/testing/testdir"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/util/cert"

	sealedsecrets "github.com/bitnami-labs/sealed-secrets/pkg/apis/sealed-secrets/v1alpha1"

	// Register Auth providers
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func TestSigning(t *testing.T) {
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
	kindCleanup, err := kind.New("test-kind", path.Join(tfdir, kubeconfigfile))
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
			"signing_cert": signingCert,
			"signing_key":  signingKey,
		},
	}

	// At the end of the test, run `terraform destroy`
	// TODO: Determine why this doesn't work
	// defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`
	result := terraform.InitAndApply(t, terraformOptions)
	t.Logf("Stdout: %v", result)

	ss, err := sealTestSecret(signingCert)
	if err != nil {
		t.Fatal(err)
	}
	ss.APIVersion = "bitnami.com/v1alpha1"
	ss.Kind = "SealedSecret"
	ssJSON, err := json.Marshal(ss)
	if err != nil {
		t.Fatal(err)
	}

	sealedSecretFile := "secret.json"
	sealedSecretPath := path.Join(tfdir, sealedSecretFile)
	err = ioutil.WriteFile(sealedSecretPath, ssJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := os.Remove(sealedSecretPath)
		if err != nil {
			t.Log(err)
		}
	}()

	// Create secret
	output, err := execKubectl(tfdir, kubeconfigfile, []string{"apply", "-f", sealedSecretFile})
	if err != nil {
		t.Fatal(err)
	}

	// Sleep for 1s
	time.Sleep(time.Second)

	// Get new secret value
	output, err = execKubectl(tfdir, kubeconfigfile, []string{"get", "secret", "mysecret", "-o", "jsonpath=\"{.data.foo}\""})
	if err != nil {
		t.Log(output)
		t.Fatal(err)
	}
	encoded := strings.Replace(string(output), "\"", "", 2)
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Error(err)
	}
	if string(data) != secretValue {
		t.Errorf("generated secret incorrect. expected %q, got %q", secretValue, string(data))
	}
}

func execKubectl(tfdir, kubeconfigfile string, params []string) (string, error) {
	cmd := exec.Command("kubectl", params...)
	cmd.Env = []string{
		fmt.Sprintf("KUBECONFIG=%s", kubeconfigfile),
	}
	cmd.Dir = tfdir
	output, err := cmd.CombinedOutput()
	return string(output), err
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
module "secrets" {
    source = "../"
    signing_cert = var.signing_cert
    signing_key = var.signing_key
}

provider "kubernetes" {
    config_path = "${path.module}/kindconfig"
}

provider "helm" {
    kubernetes {
        config_path = "${path.module}/kindconfig"
    }
}

variable signing_cert {}
variable signing_key {}
`

var secretObj = &v1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "mysecret",
		Namespace: "default",
	},
	Data: map[string][]byte{
		"foo": []byte(secretValue),
	},
}
