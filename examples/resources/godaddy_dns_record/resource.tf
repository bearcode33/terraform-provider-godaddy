# A record
resource "godaddy_dns_record" "a_record" {
  domain = "example.com"
  type   = "A"
  name   = "@"
  data   = "192.0.2.1"
  ttl    = 3600
}

# CNAME record
resource "godaddy_dns_record" "www" {
  domain = "example.com"
  type   = "CNAME"
  name   = "www"
  data   = "example.com"
  ttl    = 3600
}

# MX record
resource "godaddy_dns_record" "mx" {
  domain   = "example.com"
  type     = "MX"
  name     = "@"
  data     = "mail.example.com"
  priority = 10
  ttl      = 3600
}

# TXT record for SPF
resource "godaddy_dns_record" "spf" {
  domain = "example.com"
  type   = "TXT"
  name   = "@"
  data   = "v=spf1 include:_spf.google.com ~all"
  ttl    = 3600
}

# SRV record
resource "godaddy_dns_record" "srv" {
  domain   = "example.com"
  type     = "SRV"
  name     = "_sip._tcp"
  data     = "sipserver.example.com"
  priority = 10
  weight   = 60
  port     = 5060
  ttl      = 3600
}