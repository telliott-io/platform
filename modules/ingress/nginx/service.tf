resource "kubernetes_service" "ingress_nginx" {
  depends_on = [kubernetes_namespace.ingress_nginx]

  metadata {
    name      = "ingress-nginx"
    namespace = "ingress-nginx"

    labels = {
      "app.kubernetes.io/name" = "ingress-nginx"

      "app.kubernetes.io/part-of" = "ingress-nginx"
    }
  }

  spec {
    port {
      name        = "http"
      port        = 80
      target_port = "http"
    }

    port {
      name        = "https"
      port        = 443
      target_port = "https"
    }

    selector = {
      "app.kubernetes.io/name" = "ingress-nginx"

      "app.kubernetes.io/part-of" = "ingress-nginx"
    }

    type                    = "LoadBalancer"
    load_balancer_ip = var.load_balancer_ip
    external_traffic_policy = "Local"
  }
}

data "kubernetes_service" "ingress_nginx" {
  metadata {
    name      = "ingress-nginx"
    namespace = "ingress-nginx"
  }

  depends_on = [
    kubernetes_service.ingress_nginx,
  ]
}

output external_ip {
  value = data.kubernetes_service.ingress_nginx.load_balancer_ingress.0.ip
}