resource "kubernetes_ingress" "argocd" {
  depends_on = [helm_release.argocd]

  metadata {
    name      = "argocd-ingress"
    namespace = "argocd"
    annotations = {
      "kubernetes.io/ingress.class"                    = "nginx"
      "nginx.ingress.kubernetes.io/server-alias"       = "argocd.*"
      "nginx.ingress.kubernetes.io/force-ssl-redirect" = "true"
      "nginx.ingress.kubernetes.io/ssl-passthrough" = "true"
      "nginx.ingress.kubernetes.io/backend-protocol" = "HTTPS"
    }
  }

  spec {
    rule {
      host = "argocd"

      http {
        path {
          path = "/"

          backend {
            service_name = "argo-argocd-server"
            service_port = "https"
          }
        }
      }
    }
  }
}
