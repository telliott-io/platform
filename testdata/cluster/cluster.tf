
module "cluster" {
  source   = "github.com/telliott-io/kube-clusters//digitalocean?ref=v0.3.0"
  cluster_name = var.cluster_name
}

provider "digitalocean" {
  token = var.do_token
}

variable "do_token" {}
variable "cluster_name" {}

output "config" {
	value = module.cluster.kubernetes
	sensitive = true
}
