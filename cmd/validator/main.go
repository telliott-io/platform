package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var hostname = flag.String("hostname", "", "ingress hostname for cluster")
var ip = flag.String("ip", "", "ingress ip address for cluster")
var protocol = flag.String("protocol", "https", "protocol for request (should be http or https)")

func main() {
	flag.Parse()

	if err := verifyEnvironmentEndpoint(*hostname, *ip, *protocol); err != nil {
		log.Fatal(err)
	}

	if err := waitForArgoCD(*hostname, *ip, *protocol); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Checks succeeded")
}

func verifyEnvironmentEndpoint(hostname string, ip string, protocol string) error {
	_, err := insecureRequest(fmt.Sprintf("%v://%v/environment", protocol, ip), hostname, "HEAD")
	if err != nil {
		return err
	}
	return nil
}

type ArgoApplications struct {
	Items []struct {
		Metadata struct {
			Name string `json:"name"`
		}
		Status struct {
			Sync struct {
				Status string `json:"status"`
			} `json:"sync"`
			Health struct {
				Status string `json:"status"`
			}
		} `json:"status"`
	} `json:"items"`
}

func waitForArgoCD(hostname string, ip string, protocol string) error {
	for true {
		result, err := isArgoSyncedAndHealthy(hostname, ip, protocol)
		if err != nil {
			return err
		}
		if result {
			break
		}
		fmt.Println("Waiting for ArgoCD to sync")
		time.Sleep(time.Second)
	}
	return nil
}

func isArgoSyncedAndHealthy(hostname string, ip string, protocol string) (bool, error) {
	res, err := insecureRequest(
		fmt.Sprintf("%v://%v/api/v1/applications", protocol, ip),
		fmt.Sprintf("argocd.%v", hostname),
		"GET",
	)
	if err != nil {
		return false, err
	}
	var applications ArgoApplications
	err = json.Unmarshal([]byte(res), &applications)
	if err != nil {
		return false, err
	}
	for _, item := range applications.Items {
		if item.Status.Health.Status != "Healthy" && item.Status.Sync.Status != "Synced" {
			return false, nil
		}
	}
	return true, nil
}

func insecureRequest(url string, hostname string, method string) (string, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", err
	}
	req.Host = hostname

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	if method == "HEAD" {
		return "", nil
	}
	defer resp.Body.Close()
	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
