
module "ingress" {
  source   = "./modules/ingress/nginx"
}

module "cd" {
  source   = "./modules/argocd"
  argocd_admin_password = var.argocd_admin_password
  bootstrap_repository = var.bootstrap_repository
  bootstrap_chart = var.bootstrap_chart
  bootstrap_version = var.bootstrap_version
}

module "environment" {
  source = "./modules/envserver"
  environment_data = {
    environment = var.environment
  }
  hostname = var.hostname
}
