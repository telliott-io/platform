resource "helm_release" "argocd" {

  name       = "argo"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  namespace  = "argocd"
  create_namespace = true
  version    = "2.3.4"

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
    name = "installCRDs"
    value = false
  }
}