resource "kubernetes_namespace" "environment" {
  metadata {
    name = "environment"
  }
}