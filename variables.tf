
# Keypair for use with sealed secrets
variable secret_signing_cert {}
variable secret_signing_key {}

# Environment name for environment endpoint
variable environment {}

# Hostname for ingress
variable hostname {}

# Admin passwod for accessing argo
variable argocd_admin_password {}

# Helm chart for bootstrap ArgoCD application
variable bootstrap_repository {}
variable bootstrap_chart {}