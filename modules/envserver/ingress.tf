resource "kubernetes_ingress_v1" "envserver" {
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
            service {
              name = "envserver"
              port {
                name = "http"
              }
            }
          }
        }
      }
    }
  }
}
