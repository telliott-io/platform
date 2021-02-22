terraform {
  required_providers {
    helm = {
      source = "hashicorp/helm"
      version = "2.0.2"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = "2.0.2"
    }
  }
}

resource "kubernetes_namespace" "secrets" {
  metadata {
    name = "secrets"
  }
}

resource "helm_release" "sealed-secrets" {
  depends_on = [kubernetes_secret.signing-cert]
  name  = "sealed-secrets-controller"
  repository = "https://bitnami-labs.github.io/sealed-secrets"
  chart = "sealed-secrets"
  namespace = "secrets"
  version = "1.13.2"

  set {
    name = "secretName"
    value = "secret-signing-certs"
    type = "string"
  }
}