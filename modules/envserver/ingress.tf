resource "kubernetes_ingress" "envserver" {
  depends_on = [
    kubernetes_deployment.envserver,
  ]

  metadata {
    name      = "envserver"
    namespace = "environment"
  }

  spec {
    rule {
      host = var.hostname
      http {
        path {
          path = "/environment"

          backend {
            service_name = "envserver"
            service_port = "http"
          }
        }
      }
    }
  }
}
