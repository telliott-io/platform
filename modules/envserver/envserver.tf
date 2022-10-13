terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = "2.14.0"
    }
  }
}


resource "kubernetes_deployment" "envserver" {
  depends_on = [kubernetes_namespace.environment, kubernetes_config_map.environment]

  metadata {
    name      = "envserver"
    namespace = "environment"

    labels = {
      "app" = "envserver"
    }
  }

  spec {
    replicas = 3

    selector {
      match_labels = {
        "app" = "envserver"
      }
    }

    template {
      metadata {
        labels = {
          "app" = "envserver"
        }

        annotations = {
          "prometheus.io/scrape" = "true"
        }
      }

      spec {
        container {
          name  = "envserver"
          image = "telliottio/envserver:v0.1.0"

          port {
            name           = "http"
            container_port = 8080
          }

          env_from {
              config_map_ref {
                  name = "environment"
              }
          }

          liveness_probe {
            http_get {
              path   = "/"
              port   = "8080"
              scheme = "HTTP"
            }

            initial_delay_seconds = 10
            timeout_seconds       = 10
            period_seconds        = 10
            success_threshold     = 1
            failure_threshold     = 3
          }

          readiness_probe {
            http_get {
              path   = "/"
              port   = "8080"
              scheme = "HTTP"
            }

            timeout_seconds   = 10
            period_seconds    = 10
            success_threshold = 1
            failure_threshold = 3
          }

        }
      }
    }
  }
}

resource "kubernetes_service" "envserver" {
  depends_on = [kubernetes_namespace.environment]

  metadata {
    name      = "envserver"
    namespace = "environment"
  }

  spec {
    port {
      name        = "http"
      port        = 8080
      target_port = "http"
    }

    selector = {
      "app" = "envserver"
    }
  }
}