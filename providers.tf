
provider "kubernetes" {
    config_path = lookup(jsondecode(var.kubernetes), "config_path", null)
    host     = lookup(jsondecode(var.kubernetes), "host", null)
    username = lookup(jsondecode(var.kubernetes), "username", null)
    password = lookup(jsondecode(var.kubernetes), "password", null)
    cluster_ca_certificate = lookup(jsondecode(var.kubernetes), "cluster_ca_certificate", null)
    token = lookup(jsondecode(var.kubernetes), "token", null)
    client_certificate = lookup(jsondecode(var.kubernetes), "client_certificate", null)
    client_key = lookup(jsondecode(var.kubernetes), "client_key", null)
}

provider "helm" {
  kubernetes {
    config_path = lookup(jsondecode(var.kubernetes), "config_path", null)
    host     = lookup(jsondecode(var.kubernetes), "host", null)
    username = lookup(jsondecode(var.kubernetes), "username", null)
    password = lookup(jsondecode(var.kubernetes), "password", null)
    cluster_ca_certificate = lookup(jsondecode(var.kubernetes), "cluster_ca_certificate", null)
    token = lookup(jsondecode(var.kubernetes), "token", null)
    client_certificate = lookup(jsondecode(var.kubernetes), "client_certificate", null)
    client_key = lookup(jsondecode(var.kubernetes), "client_key", null)
  }
}

variable kubernetes {}