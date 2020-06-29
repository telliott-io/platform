resource "cloudflare_record" "ingress" {
  zone_id = var.cloudflare_zone_id
  name    = "ingress"
  value   = var.ingress_ip
  type    = "A"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "jaeger" {
  zone_id = var.cloudflare_zone_id
  name    = "jaeger"
  value   = "ingress.${var.domain}"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "prometheus" {
  zone_id = var.cloudflare_zone_id
  name    = "prometheus"
  value   = "ingress.${var.domain}"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "grafana" {
  zone_id = var.cloudflare_zone_id
  name    = "grafana"
  value   = "ingress.${var.domain}"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "emojicode" {
  zone_id = var.cloudflare_zone_id
  name    = "emojicode"
  value   = "ingress.${var.domain}"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "argocd" {
  zone_id = var.cloudflare_zone_id
  name    = "argocd"
  value   = "ingress.${var.domain}"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "root" {
  zone_id = var.cloudflare_zone_id
  name    = "@"
  value   = "ingress.${var.domain}"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_record" "www" {
  zone_id = var.cloudflare_zone_id
  name    = "www"
  value   = "ingress.${var.domain}"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}