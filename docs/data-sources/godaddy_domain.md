# godaddy_domain (Data Source)

Retrieves information about a GoDaddy domain including status, expiration, nameservers, and settings.

## Example Usage

### Basic Domain Information

```terraform
data "godaddy_domain" "example" {
  domain = "example.com"
}

output "domain_status" {
  value = data.godaddy_domain.example.status
}

output "expires_on" {
  value = data.godaddy_domain.example.expires
}
```

### Complete Domain Information with Outputs

```terraform
data "godaddy_domain" "example" {
  domain = "example.com"
}

# Output key domain information
output "domain_info" {
  value = {
    domain                = data.godaddy_domain.example.domain
    status                = data.godaddy_domain.example.status
    expires               = data.godaddy_domain.example.expires
    locked                = data.godaddy_domain.example.locked
    privacy               = data.godaddy_domain.example.privacy
    renew_auto            = data.godaddy_domain.example.renew_auto
    expiration_protected  = data.godaddy_domain.example.expiration_protected
    transfer_protected    = data.godaddy_domain.example.transfer_protected
    nameservers           = data.godaddy_domain.example.nameservers
    created_at            = data.godaddy_domain.example.created_at
  }
}

# Check if domain expires soon (within 30 days)
locals {
  expires_date = formatdate("YYYY-MM-DD", data.godaddy_domain.example.expires)
  expires_soon = timecmp(data.godaddy_domain.example.expires, timeadd(timestamp(), "720h")) < 0
}

output "domain_expires_soon" {
  value = local.expires_soon
}
```

### Using Domain Information in Resources

```terraform
data "godaddy_domain" "main" {
  domain = "example.com"
}

# Create DNS records only if domain is active
resource "godaddy_dns_record" "www" {
  count = data.godaddy_domain.main.status == "ACTIVE" ? 1 : 0
  
  domain = data.godaddy_domain.main.domain
  type   = "CNAME"
  name   = "www"
  data   = data.godaddy_domain.main.domain
  ttl    = 3600
}

# Use existing nameservers for subdomain delegation
resource "godaddy_dns_record" "subdomain_ns" {
  domain = data.godaddy_domain.main.domain
  type   = "NS"
  name   = "subdomain"
  data   = data.godaddy_domain.main.nameservers[0]
  ttl    = 3600
}
```

## Schema

### Required

- `domain` (String) - The domain name to retrieve information for.

### Read-Only

- `status` (String) - Current status of the domain (e.g., "ACTIVE", "EXPIRED", "PENDING").
- `expires` (String) - Domain expiration date in RFC3339 format.
- `expiration_protected` (Boolean) - Whether the domain is protected from expiration.
- `hold_registrar` (Boolean) - Whether the domain has a registrar hold.
- `locked` (Boolean) - Whether the domain is locked to prevent unauthorized transfers.
- `privacy` (Boolean) - Whether WHOIS privacy protection is enabled.
- `renew_auto` (Boolean) - Whether the domain is set to auto-renew.
- `transfer_protected` (Boolean) - Whether the domain is protected from transfers.
- `nameservers` (List of String) - Current nameservers for the domain.
- `created_at` (String) - Domain creation date in RFC3339 format.
- `modified_at` (String) - Last modification date in RFC3339 format.

## Domain Status Values

The `status` field can have the following values:

- **ACTIVE** - Domain is active and functioning normally
- **EXPIRED** - Domain has expired but may still be renewable
- **PENDING_DELETE** - Domain is pending deletion
- **PENDING_TRANSFER** - Domain transfer is in progress
- **REGISTRY_HOLD** - Domain is on registry hold
- **REGISTRAR_HOLD** - Domain is on registrar hold
- **SUSPENDED** - Domain is suspended
- **PENDING** - Domain registration is pending

## Working with Dates

The `expires`, `created_at`, and `modified_at` fields are returned in RFC3339 format. You can use Terraform's time functions to work with these dates:

```terraform
data "godaddy_domain" "example" {
  domain = "example.com"
}

locals {
  # Convert to readable format
  expires_readable = formatdate("DD MMM YYYY", data.godaddy_domain.example.expires)
  
  # Calculate days until expiration
  days_until_expiry = (parseint(substr(data.godaddy_domain.example.expires, 0, 4), 10) * 365 + 
                      parseint(substr(data.godaddy_domain.example.expires, 5, 2), 10) * 30 + 
                      parseint(substr(data.godaddy_domain.example.expires, 8, 2), 10)) - 
                     (parseint(formatdate("YYYY", timestamp()), 10) * 365 + 
                      parseint(formatdate("MM", timestamp()), 10) * 30 + 
                      parseint(formatdate("DD", timestamp()), 10))
  
  # Check if domain expires within 30 days
  expires_soon = timecmp(data.godaddy_domain.example.expires, timeadd(timestamp(), "720h")) < 0
}

output "expiration_info" {
  value = {
    expires_on    = local.expires_readable
    expires_soon  = local.expires_soon
  }
}
```

## Example Use Cases

### Monitoring Domain Expiration

```terraform
data "godaddy_domain" "domains" {
  for_each = toset(var.domains_to_monitor)
  domain   = each.value
}

# Create alerts for domains expiring soon
resource "local_file" "expiration_report" {
  filename = "domain_expiration_report.json"
  content = jsonencode({
    for domain, info in data.godaddy_domain.domains : domain => {
      expires     = info.expires
      expires_soon = timecmp(info.expires, timeadd(timestamp(), "720h")) < 0
      auto_renew  = info.renew_auto
      status      = info.status
    }
  })
}
```

### Conditional Resource Creation

```terraform
data "godaddy_domain" "example" {
  domain = "example.com"
}

# Only create records if domain is active and not expired
resource "godaddy_dns_record" "web_records" {
  for_each = data.godaddy_domain.example.status == "ACTIVE" ? var.web_records : {}
  
  domain = data.godaddy_domain.example.domain
  type   = each.value.type
  name   = each.value.name
  data   = each.value.data
  ttl    = each.value.ttl
}
```

### Backup Current Configuration

```terraform
data "godaddy_domain" "backup" {
  domain = "example.com"
}

# Store current configuration for backup/restore
output "domain_backup" {
  value = {
    domain              = data.godaddy_domain.backup.domain
    locked              = data.godaddy_domain.backup.locked
    privacy             = data.godaddy_domain.backup.privacy
    renew_auto          = data.godaddy_domain.backup.renew_auto
    expiration_protected = data.godaddy_domain.backup.expiration_protected
    nameservers         = data.godaddy_domain.backup.nameservers
  }
  sensitive = false
}
```

## Notes

### API Rate Limits

Reading domain information counts against GoDaddy's API rate limits. If you need to query multiple domains frequently, consider caching the results or using fewer queries.

### Data Freshness

Domain information is retrieved in real-time from the GoDaddy API. However, some changes (especially status changes) may take time to reflect in the API.

### Permissions

Your API credentials must have access to the domain you're querying. Domains not in your account or not accessible with your credentials will result in an error.

### Error Handling

If a domain is not found or not accessible, the data source will fail with an appropriate error message. Use this behavior to validate domain ownership in your configurations.

```terraform
# This will fail if the domain is not in your account
data "godaddy_domain" "verify_ownership" {
  domain = var.domain_to_verify
}

# Use the data source to ensure the domain exists before creating resources
resource "godaddy_dns_record" "validated_record" {
  domain = data.godaddy_domain.verify_ownership.domain
  type   = "A"
  name   = "@"
  data   = "192.0.2.1"
  ttl    = 3600
}
```