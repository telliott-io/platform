resource "kubernetes_config_map" "environment" {
  depends_on = [kubernetes_namespace.environment]
  metadata {
    name = "environment"
    namespace = "environment"
  }

  data = var.environment_data
}