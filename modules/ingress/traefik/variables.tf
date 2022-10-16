variable load_balancer_ip {
    default = null
}

variable service_type {
    type = string
    default = "LoadBalancer"
    description = "Kubernetes service type as per https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types"
}

