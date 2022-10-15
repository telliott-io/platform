terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
      version = "2.7.1"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = "2.14.0"
    }
  }
}

resource "helm_release" "traefik" {
  name       = "traefikingress"
  repository = "https://helm.traefik.io/traefik"
  chart      = "traefik"
#   namespace  = "traefik"
#   create_namespace = true
  version    = "v15.1.0"
}