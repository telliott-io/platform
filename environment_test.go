package platform

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

type EnvironmentData struct {
	Environment string `json:"environment"`
}

// checkEnvironmentHost verifies that the environment ingress returns the expected
// environment name.
func checkEnvironmentHost(t *testing.T, host string, hostHeader string, expectedPlatform string) {
	url := fmt.Sprintf("%v/environment", host)
	t.Logf("Checking environment output at %q with host header %q", url, hostHeader)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = hostHeader
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var env EnvironmentData
	if err := json.Unmarshal(body, &env); err != nil {
		t.Errorf("could not parse environment endpoint response: %v\nResponse was:\n%v", err, string(body))
	}

	if env.Environment != expectedPlatform {
		t.Errorf("expected environment to be %q, got %q", expectedPlatform, env.Environment)
	}
}
