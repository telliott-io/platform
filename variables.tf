variable environment {
    type = string
    description = "Environment name, output by the /environment endpoint as a way to uniquely identify this cluster."
}

variable hostname {
    type = string
    description = "Base hostname for ingresses."
}

variable argocd_admin_password {
    type = string
    description = "Admin password for accessing ArgoCD."
}

variable bootstrap_repository {
    type = string
    description = "Helm repository for application bootstrapped into ArgoCD."
}

variable bootstrap_chart {
    type = string
    description = "Helm chart for application bootstrapped into ArgoCD."    
}

variable bootstrap_version {
    type = string
    description = "Helm chart version for application bootstrapped into ArgoCD."      
}