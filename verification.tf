resource "null_resource" "verification" {
  depends_on = [
      module.ingress,
      module.cd,
      module.secrets,
      module.environment
  ]
  provisioner "local-exec" {
    command = "go run github.com/telliott-io/platform/cmd/validator --hostname ${var.hostname} --ip ${local.public_url}"
  }
}