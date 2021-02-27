resource "helm_release" "bootstrap" {
  depends_on = [helm_release.argocd]

  name       = "bootstrap"
  repository = var.bootstrap_repository
  chart      = var.bootstrap_chart
  namespace  = "argocd"
  version    = var.bootstrap_version
}