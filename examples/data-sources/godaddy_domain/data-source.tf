data "godaddy_domain" "example" {
  domain = "example.com"
}

output "domain_status" {
  value = data.godaddy_domain.example.status
}

output "domain_expires" {
  value = data.godaddy_domain.example.expires
}

output "domain_nameservers" {
  value = data.godaddy_domain.example.nameservers
}