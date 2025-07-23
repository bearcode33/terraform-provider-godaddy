# godaddy_dns_records (Data Source)

Retrieves DNS records for a GoDaddy domain with optional filtering by record type and name.

## Example Usage

### Get All DNS Records

```terraform
data "godaddy_dns_records" "all" {
  domain = "example.com"
}

output "all_records" {
  value = data.godaddy_dns_records.all.records
}
```

### Filter by Record Type

```terraform
# Get all A records
data "godaddy_dns_records" "a_records" {
  domain = "example.com"
  type   = "A"
}

# Get all MX records
data "godaddy_dns_records" "mx_records" {
  domain = "example.com"
  type   = "MX"
}

output "mail_servers" {
  value = [for record in data.godaddy_dns_records.mx_records.records : {
    server   = record.data
    priority = record.priority
  }]
}
```

### Filter by Record Name

```terraform
# Get all records for www subdomain
data "godaddy_dns_records" "www_records" {
  domain = "example.com"
  name   = "www"
}

# Get all records for root domain
data "godaddy_dns_records" "root_records" {
  domain = "example.com"
  name   = "@"
}
```

### Filter by Both Type and Name

```terraform
# Get specific A record for www
data "godaddy_dns_records" "www_a" {
  domain = "example.com"
  type   = "A"
  name   = "www"
}

# Get MX records for root domain
data "godaddy_dns_records" "root_mx" {
  domain = "example.com"
  type   = "MX"
  name   = "@"
}
```

### Complex Filtering and Processing

```terraform
data "godaddy_dns_records" "all" {
  domain = "example.com"
}

locals {
  # Group records by type
  records_by_type = {
    for record in data.godaddy_dns_records.all.records :
    record.type => record...
  }
  
  # Extract all IP addresses from A records
  ip_addresses = [
    for record in data.godaddy_dns_records.all.records :
    record.data if record.type == "A"
  ]
  
  # Find records with low TTL (less than 1 hour)
  low_ttl_records = [
    for record in data.godaddy_dns_records.all.records :
    record if record.ttl < 3600
  ]
  
  # Extract mail servers with priorities
  mail_servers = [
    for record in data.godaddy_dns_records.all.records : {
      server   = record.data
      priority = record.priority
    } if record.type == "MX"
  ]
}

output "dns_analysis" {
  value = {
    total_records     = length(data.godaddy_dns_records.all.records)
    records_by_type   = { for k, v in local.records_by_type : k => length(v) }
    ip_addresses      = local.ip_addresses
    low_ttl_count     = length(local.low_ttl_records)
    mail_servers      = local.mail_servers
  }
}
```

### Using Records in Resource Creation

```terraform
data "godaddy_dns_records" "existing_a" {
  domain = "example.com"
  type   = "A"
}

# Create AAAA records for existing A records
resource "godaddy_dns_record" "ipv6" {
  for_each = {
    for record in data.godaddy_dns_records.existing_a.records :
    record.name => record
  }
  
  domain = "example.com"
  type   = "AAAA"
  name   = each.value.name
  data   = "2001:db8::${index(data.godaddy_dns_records.existing_a.records, each.value) + 1}"
  ttl    = each.value.ttl
}
```

### Backup DNS Configuration

```terraform
data "godaddy_dns_records" "backup" {
  domain = var.domain_to_backup
}

# Export DNS records to JSON file
resource "local_file" "dns_backup" {
  filename = "${var.domain_to_backup}_dns_backup.json"
  content = jsonencode({
    domain = var.domain_to_backup
    backup_date = timestamp()
    records = data.godaddy_dns_records.backup.records
  })
}

# Create outputs for easy viewing
output "dns_backup_summary" {
  value = {
    domain        = var.domain_to_backup
    total_records = length(data.godaddy_dns_records.backup.records)
    record_types  = distinct([for r in data.godaddy_dns_records.backup.records : r.type])
    backup_file   = local_file.dns_backup.filename
  }
}
```

## Schema

### Required

- `domain` (String) - The domain name to retrieve DNS records for.

### Optional

- `type` (String) - Filter records by type (A, AAAA, CNAME, MX, TXT, SRV, etc.).
- `name` (String) - Filter records by name/subdomain.

### Read-Only

- `records` (List of Object) - List of DNS records matching the filter criteria.

### Records Object Schema

Each record in the `records` list contains:

- `type` (String) - The DNS record type.
- `name` (String) - The record name/subdomain.
- `data` (String) - The record value.
- `ttl` (Number) - Time to live in seconds.
- `priority` (Number) - Priority for MX, SRV, NAPTR, URI records (null if not applicable).
- `port` (Number) - Port for SRV records (null if not applicable).
- `weight` (Number) - Weight for SRV, URI records (null if not applicable).
- `service` (String) - Service for SRV records (null if not applicable).
- `protocol` (String) - Protocol for SRV records (null if not applicable).

## Filtering Behavior

### No Filters
Returns all DNS records for the domain.

### Type Filter Only
Returns all records of the specified type, regardless of name.

### Name Filter Only
Returns all records with the specified name, regardless of type.

### Both Type and Name Filters
Returns only records that match both the specified type and name.

## Working with Record Data

### Processing Different Record Types

```terraform
data "godaddy_dns_records" "all" {
  domain = "example.com"
}

locals {
  # Process A records
  a_records = [
    for record in data.godaddy_dns_records.all.records :
    {
      name = record.name
      ip   = record.data
      ttl  = record.ttl
    } if record.type == "A"
  ]
  
  # Process MX records with priorities
  mx_records = [
    for record in data.godaddy_dns_records.all.records :
    {
      name     = record.name
      server   = record.data
      priority = record.priority
      ttl      = record.ttl
    } if record.type == "MX"
  ]
  
  # Process SRV records with all fields
  srv_records = [
    for record in data.godaddy_dns_records.all.records :
    {
      name     = record.name
      target   = record.data
      priority = record.priority
      weight   = record.weight
      port     = record.port
      service  = record.service
      protocol = record.protocol
      ttl      = record.ttl
    } if record.type == "SRV"
  ]
  
  # Process TXT records for specific purposes
  spf_records = [
    for record in data.godaddy_dns_records.all.records :
    record.data if record.type == "TXT" && startswith(record.data, "v=spf1")
  ]
  
  dmarc_records = [
    for record in data.godaddy_dns_records.all.records :
    record.data if record.type == "TXT" && startswith(record.data, "v=DMARC1")
  ]
}
```

### DNS Health Checks

```terraform
data "godaddy_dns_records" "health_check" {
  domain = "example.com"
}

locals {
  # Check for common DNS issues
  dns_issues = {
    # Records with very short TTL (< 5 minutes)
    short_ttl_records = [
      for record in data.godaddy_dns_records.health_check.records :
      record if record.ttl < 300
    ]
    
    # Multiple A records for root (potential issue)
    multiple_root_a = length([
      for record in data.godaddy_dns_records.health_check.records :
      record if record.type == "A" && record.name == "@"
    ]) > 1
    
    # Missing common records
    missing_www = length([
      for record in data.godaddy_dns_records.health_check.records :
      record if record.name == "www"
    ]) == 0
    
    missing_mx = length([
      for record in data.godaddy_dns_records.health_check.records :
      record if record.type == "MX"
    ]) == 0
    
    # Check for SPF record
    has_spf = length([
      for record in data.godaddy_dns_records.health_check.records :
      record if record.type == "TXT" && contains(record.data, "v=spf1")
    ]) > 0
  }
}

output "dns_health_report" {
  value = local.dns_issues
}
```

## Example Use Cases

### DNS Migration Planning

```terraform
# Get current DNS configuration from source
data "godaddy_dns_records" "source" {
  domain = "example.com"
}

# Create equivalent records on new provider
resource "other_provider_dns_record" "migrated" {
  for_each = {
    for record in data.godaddy_dns_records.source.records :
    "${record.type}-${record.name}" => record
  }
  
  domain = "example.com"
  type   = each.value.type
  name   = each.value.name
  value  = each.value.data
  ttl    = each.value.ttl
  
  # Handle type-specific fields
  priority = each.value.priority
  weight   = each.value.weight
  port     = each.value.port
}
```

### Load Balancer Configuration

```terraform
data "godaddy_dns_records" "web_servers" {
  domain = "example.com"
  type   = "A"
  name   = "api"
}

# Use DNS records to configure load balancer
resource "aws_lb_target_group_attachment" "web_servers" {
  for_each = {
    for record in data.godaddy_dns_records.web_servers.records :
    record.data => record
  }
  
  target_group_arn = aws_lb_target_group.web.arn
  target_id        = each.key
  port             = 80
}
```

### DNS Validation

```terraform
data "godaddy_dns_records" "validation" {
  domain = "example.com"
}

# Validate required DNS records exist
check "required_dns_records" {
  assert {
    condition = length([
      for record in data.godaddy_dns_records.validation.records :
      record if record.type == "MX"
    ]) > 0
    error_message = "Domain must have at least one MX record configured."
  }
  
  assert {
    condition = length([
      for record in data.godaddy_dns_records.validation.records :
      record if record.type == "TXT" && contains(record.data, "v=spf1")
    ]) > 0
    error_message = "Domain must have an SPF record configured."
  }
}
```

## Notes

### Performance Considerations

- Retrieving all records for a domain with many DNS records may take time
- Use specific filters (`type` and `name`) when possible to reduce API calls
- Consider caching results for frequently accessed data

### API Rate Limits

DNS record queries count against GoDaddy's API rate limits. Be mindful of this when querying multiple domains or making frequent requests.

### Record Ordering

Records are returned in the order provided by the GoDaddy API. Don't rely on a specific ordering for your logic.

### Null Values

Fields that don't apply to specific record types will be `null`. For example, `priority` will be `null` for A records, and `port` will be `null` for non-SRV records.

### Large Datasets

For domains with many DNS records, consider using pagination or specific filtering to avoid timeouts and reduce memory usage.