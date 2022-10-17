resource "kubernetes_ingress_v1" "argocd" {
  depends_on = [helm_release.argocd]

  metadata {
    name      = "argocd-ingress"
    namespace = "argocd"
  }

  spec {
    rule {
      host = "argocd.${var.hostname}"

      http {
        path {
          path = "/"

          backend {
            service {
              name = "argo-argocd-server"
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
