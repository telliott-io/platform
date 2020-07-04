locals {
    public_url = coalesce(var.dns_name, module.ingress.external_ip)
}