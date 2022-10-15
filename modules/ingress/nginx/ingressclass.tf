resource "kubernetes_ingress_class_v1" "nginx-ingress-class" {
  metadata {
    name = "nginx"
  }

  spec {
    controller = "nginx.org/ingress-controller"
  }
}