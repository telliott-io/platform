variable bootstrap_repository {
    default = null
}
variable bootstrap_chart {
    default = null
}

variable bootstrap_version {
    // Previous version was hard-coded, add default for backwards compatibility
    default = "0.1.3"
}

variable argocd_admin_password {}

variable hostname {}