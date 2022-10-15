resource "kubernetes_ingress_class_v1" "nginx-ingress-class" {
  metadata {
    name = "nginx"
    annotations = {
      "ingressclass.kubernetes.io/is-default-class" = "true"
    }
  }

  spec {
    controller = "nginx.org/ingress-controller"
    
  }
}