# Get all DNS records for a domain
data "godaddy_dns_records" "all" {
  domain = "example.com"
}

# Get all A records
data "godaddy_dns_records" "a_records" {
  domain = "example.com"
  type   = "A"
}

# Get all records for www subdomain
data "godaddy_dns_records" "www_records" {
  domain = "example.com"
  name   = "www"
}

# Get specific record
data "godaddy_dns_records" "mx_records" {
  domain = "example.com"
  type   = "MX"
  name   = "@"
}

output "all_records" {
  value = data.godaddy_dns_records.all.records
}

output "a_record_ips" {
  value = [for r in data.godaddy_dns_records.a_records.records : r.data]
}