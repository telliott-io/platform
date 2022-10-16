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
  namespace  = "traefik"
  create_namespace = true
  version    = "v15.1.0"

  dynamic "set" {
    for_each = (var.load_balancer_ip != null) ? [var.load_balancer_ip] : []
    content {
        name = "service.externalIps[0]"
        value = set.value
    }
  }

  set {
    name = "service.type"
    value = var.service_type
  }

  dynamic "set" {
    for_each = (var.service_type == "NodePort") ? ["32080"] : []
    content {
      name = "ports.web.nodePort"
      value = set.value
    }
  }
}