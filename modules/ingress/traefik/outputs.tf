data "kubernetes_service" "ingress_traefik" {
  depends_on = [helm_release.traefik]
  metadata {
    name      = "traefikingress"
    namespace = "traefik"
  }
}

output external_ip {
  value = data.kubernetes_service.ingress_traefik.status != null ? data.kubernetes_service.ingress_traefik.status.0.load_balancer.0.ingress.0.ip : null
}