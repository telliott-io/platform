output "ingress_address" {
    depends_on = [null_resource.verification]
    value = "${module.ingress.external_ip}"
}