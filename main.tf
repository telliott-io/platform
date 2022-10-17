
module "ingress" {
  source   = "./modules/ingress/traefik"
  load_balancer_ip = var.load_balancer_ip
  service_type = var.service_type
}

module "cd" {
  source   = "./modules/argocd"
  argocd_admin_password = var.argocd_admin_password
  bootstrap_repository = var.bootstrap_repository
  bootstrap_chart = var.bootstrap_chart
  bootstrap_version = var.bootstrap_version
  
  hostname = var.hostname
}

module "environment" {
  source = "./modules/envserver"
  environment_data = {
    environment = var.environment
  }
  hostname = var.hostname
}
