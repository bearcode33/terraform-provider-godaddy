terraform {
  required_providers {
    godaddy = {
      source  = "bearcode33/godaddy"
      version = "~> 1.0"
    }
  }
}

provider "godaddy" {
  # Configuration via environment variables
}

variable "domain" {
  description = "Domain name to manage"
  type        = string
  default     = "example.com"
}

# A Record - IPv4 address
resource "godaddy_dns_record" "a_record" {
  domain = var.domain
  type   = "A"
  name   = "@"
  data   = "192.0.2.1"
  ttl    = 3600
}

# AAAA Record - IPv6 address
resource "godaddy_dns_record" "aaaa_record" {
  domain = var.domain
  type   = "AAAA"
  name   = "@"
  data   = "2001:db8::1"
  ttl    = 3600
}

# CNAME Record - Canonical name
resource "godaddy_dns_record" "cname_www" {
  domain = var.domain
  type   = "CNAME"
  name   = "www"
  data   = "example.com"
  ttl    = 3600
}

# MX Record - Mail exchange
resource "godaddy_dns_record" "mx_primary" {
  domain   = var.domain
  type     = "MX"
  name     = "@"
  data     = "mail.example.com"
  priority = 10
  ttl      = 3600
}

resource "godaddy_dns_record" "mx_secondary" {
  domain   = var.domain
  type     = "MX"
  name     = "@"
  data     = "mail2.example.com"
  priority = 20
  ttl      = 3600
}

# TXT Record - Text record
resource "godaddy_dns_record" "txt_spf" {
  domain = var.domain
  type   = "TXT"
  name   = "@"
  data   = "v=spf1 include:_spf.google.com ~all"
  ttl    = 3600
}

# TXT Record - DKIM
resource "godaddy_dns_record" "txt_dkim" {
  domain = var.domain
  type   = "TXT"
  name   = "google._domainkey"
  data   = "v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC..."
  ttl    = 3600
}

# TXT Record - DMARC
resource "godaddy_dns_record" "txt_dmarc" {
  domain = var.domain
  type   = "TXT"
  name   = "_dmarc"
  data   = "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.com"
  ttl    = 3600
}

# SRV Record - Service record
resource "godaddy_dns_record" "srv_sip_tcp" {
  domain   = var.domain
  type     = "SRV"
  name     = "_sip._tcp"
  data     = "sipserver.example.com"
  priority = 10
  weight   = 60
  port     = 5060
  service  = "_sip"
  protocol = "_tcp"
  ttl      = 3600
}

# SRV Record - XMPP
resource "godaddy_dns_record" "srv_xmpp_tcp" {
  domain   = var.domain
  type     = "SRV"
  name     = "_xmpp-server._tcp"
  data     = "xmpp.example.com"
  priority = 5
  weight   = 0
  port     = 5269
  service  = "_xmpp-server"
  protocol = "_tcp"
  ttl      = 3600
}

# CAA Record - Certificate Authority Authorization
resource "godaddy_dns_record" "caa_letsencrypt" {
  domain = var.domain
  type   = "CAA"
  name   = "@"
  data   = "0 issue \"letsencrypt.org\""
  ttl    = 3600
}

resource "godaddy_dns_record" "caa_wildcard" {
  domain = var.domain
  type   = "CAA"
  name   = "@"
  data   = "0 issuewild \"letsencrypt.org\""
  ttl    = 3600
}

# NS Record - Name server
resource "godaddy_dns_record" "ns_subdomain" {
  domain = var.domain
  type   = "NS"
  name   = "subdomain"
  data   = "ns1.subdomain.example.com"
  ttl    = 3600
}

# PTR Record - Pointer record (for reverse DNS)
resource "godaddy_dns_record" "ptr_reverse" {
  domain = var.domain
  type   = "PTR"
  name   = "1.2.0.192.in-addr.arpa"
  data   = "server.example.com"
  ttl    = 3600
}

# SSHFP Record - SSH fingerprint
resource "godaddy_dns_record" "sshfp" {
  domain = var.domain
  type   = "SSHFP"
  name   = "@"
  data   = "1 1 123456789abcdef67890123456789abcdef67890"
  ttl    = 3600
}

# TLSA Record - TLS certificate association
resource "godaddy_dns_record" "tlsa" {
  domain = var.domain
  type   = "TLSA"
  name   = "_443._tcp"
  data   = "3 1 1 0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
  ttl    = 3600
}

# DS Record - Delegation Signer
resource "godaddy_dns_record" "ds" {
  domain = var.domain
  type   = "DS"
  name   = "subdomain"
  data   = "12345 7 1 1234567890ABCDEF1234567890ABCDEF12345678"
  ttl    = 3600
}

# NAPTR Record - Naming Authority Pointer
resource "godaddy_dns_record" "naptr" {
  domain   = var.domain
  type     = "NAPTR"
  name     = "@"
  data     = "100 10 \"u\" \"E2U+sip\" \"!^.*$!sip:info@example.com!\" ."
  priority = 100
  ttl      = 3600
}

# URI Record - Uniform Resource Identifier
resource "godaddy_dns_record" "uri" {
  domain   = var.domain
  type     = "URI"
  name     = "_http._tcp"
  data     = "https://example.com/path"
  priority = 10
  weight   = 1
  ttl      = 3600
}

# LOC Record - Location
resource "godaddy_dns_record" "loc" {
  domain = var.domain
  type   = "LOC"
  name   = "@"
  data   = "37 23 30.000 N 122 02 30.000 W 10m 100m 100m 10m"
  ttl    = 3600
}

# CERT Record - Certificate
resource "godaddy_dns_record" "cert" {
  domain = var.domain
  type   = "CERT"
  name   = "@"
  data   = "1 0 0 MIGfMA0GCSq..."
  ttl    = 3600
}

# DNAME Record - Delegation Name
resource "godaddy_dns_record" "dname" {
  domain = var.domain
  type   = "DNAME"
  name   = "old"
  data   = "new.example.com"
  ttl    = 3600
}