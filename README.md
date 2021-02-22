# platform

A Terraform module for provisioning standard resources on a Kubernetes cluster.

## What is deployed

The parent module is composed of several sub modules.

- argocd: A deployment of [ArgoCD](https://argoproj.github.io/argo-cd/), using a bootstrap application to load other software.
- envserver: A test server providing environment information on an external endpoint for verification.
- ingress: A standard ingress controller with associated Load Balancer.
- secrets: A deployment of [sealed-secrets](https://github.com/bitnami-labs/sealed-secrets) to assist with managing secrets that may be stored, encrypted, in public repos.

A combination of Terraform's Kubernetes and Helm providers is used to deploy these modules.

## Testing

Tests are configured with GitHub Actions, and may be run locally using [act](https://github.com/nektos/act). Using the following command at the root of this repo to test all modules and their integration:

    act

All underlying tests are written using [Terratest](https://github.com/gruntwork-io/terratest) in Go test files.
