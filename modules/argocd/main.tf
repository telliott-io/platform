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

resource "helm_release" "argocd" {
  name       = "argo"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  namespace  = "argocd"
  create_namespace = true
  version    = "5.5.24"

  set_sensitive {
    name  = "configs.secret.argocdServerAdminPassword"
    value = var.argocd_admin_password
  }

  set {
    name = "server.rbacConfig.policy\\.default"
    value = "role:readonly"
  }

  set {
    name = "server.config.users\\.anonymous\\.enabled"
    value = "true"
    type = "string"
  }

  set {
    name = "configs.params.server\\.insecure"
    value = true
  }
}