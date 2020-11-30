resource "helm_release" "tictactoe" {
  name       = "tictactoe"
  repository = "https://theothertomelliott.github.io/tic-tac-toverengineered/"
  chart      = "tic-tac-toe"
  namespace  = "tictactoe"
  version    = "v0.1.27"
  create_namespace = true

  set {
    name = "hostname"
    value = "tictactoe.${var.hostname}"
    type = "string"
  }

}