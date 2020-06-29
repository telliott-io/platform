resource "kubernetes_namespace" "secrets" {
  metadata {
    name = "secrets"
  }
}

resource "helm_release" "sealed-secrets" {
  depends_on = [kubernetes_secret.signing-cert]
  name  = "sealed-secrets-controller"
  repository = "https://kubernetes-charts.storage.googleapis.com"
  chart = "sealed-secrets"
  namespace = "secrets"
  version = "1.10.0"

  set {
    name = "secretName"
    value = "secret-signing-certs"
    type = "string"
  }
}