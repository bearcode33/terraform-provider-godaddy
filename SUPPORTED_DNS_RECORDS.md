# DNS Record Types Support

This document details all DNS record types supported by the GoDaddy Terraform Provider.

## ✅ Fully Supported Record Types

### Basic Records
| Type | Description | Required Fields | Optional Fields |
|------|-------------|----------------|-----------------|
| **A** | IPv4 address | `data` | `ttl` |
| **AAAA** | IPv6 address | `data` | `ttl` |
| **CNAME** | Canonical name (alias) | `data` | `ttl` |
| **TXT** | Text record | `data` | `ttl` |
| **NS** | Name server | `data` | `ttl` |
| **PTR** | Pointer (reverse DNS) | `data` | `ttl` |

### Mail Records
| Type | Description | Required Fields | Optional Fields |
|------|-------------|----------------|-----------------|
| **MX** | Mail exchange | `data`, `priority` | `ttl` |

### Service Records
| Type | Description | Required Fields | Optional Fields |
|------|-------------|----------------|-----------------|
| **SRV** | Service location | `data`, `priority`, `weight`, `port`, `service`, `protocol` | `ttl` |

### Security Records
| Type | Description | Required Fields | Optional Fields |
|------|-------------|----------------|-----------------|
| **CAA** | Certificate Authority Authorization | `data` | `ttl` |
| **SSHFP** | SSH fingerprint | `data` | `ttl` |
| **TLSA** | TLS certificate association | `data` | `ttl` |
| **DS** | DNSSEC Delegation Signer | `data` | `ttl` |

### Advanced Records
| Type | Description | Required Fields | Optional Fields |
|------|-------------|----------------|-----------------|
| **NAPTR** | Naming Authority Pointer | `data`, `priority` | `ttl` |
| **URI** | Uniform Resource Identifier | `data`, `priority`, `weight` | `ttl` |
| **LOC** | Geographic location | `data` | `ttl` |
| **CERT** | Certificate | `data` | `ttl` |
| **DNAME** | Delegation Name | `data` | `ttl` |
| **SOA** | Start of Authority | `data` | `ttl` |

## Field Specifications

### Common Fields
- **domain**: Domain name (required for all records)
- **type**: Record type (required for all records)
- **name**: Record name/subdomain (required for all records)
- **data**: Record value (required for all records)
- **ttl**: Time to live in seconds (optional, default: 3600)

### Type-Specific Fields
- **priority**: Used for MX, SRV, NAPTR, URI records (0-65535)
- **weight**: Used for SRV, URI records (0-65535)
- **port**: Used for SRV records (1-65535)
- **service**: Used for SRV records (must start with _)
- **protocol**: Used for SRV records (must start with _)

## Data Format Examples

### A Record
```hcl
data = "192.0.2.1"
```

### AAAA Record
```hcl
data = "2001:db8::1"
```

### MX Record
```hcl
data = "mail.example.com"
priority = 10
```

### TXT Record
```hcl
# SPF
data = "v=spf1 include:_spf.google.com ~all"

# DKIM
data = "v=DKIM1; k=rsa; p=MIGfMA0GCSq..."

# DMARC
data = "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.com"
```

### SRV Record
```hcl
data = "target.example.com"
priority = 10
weight = 60
port = 5060
service = "_sip"
protocol = "_tcp"
```

### CAA Record
```hcl
# Issue authorization
data = "0 issue \"letsencrypt.org\""

# Wildcard authorization
data = "0 issuewild \"letsencrypt.org\""

# Incident reporting
data = "0 iodef \"mailto:admin@example.com\""
```

### SSHFP Record
```hcl
data = "1 1 123456789abcdef67890123456789abcdef67890"
# Format: algorithm fingerprint_type fingerprint
```

### TLSA Record
```hcl
data = "3 1 1 0123456789abcdef..."
# Format: cert_usage selector matching_type cert_data
```

### DS Record
```hcl
data = "12345 7 1 1234567890ABCDEF..."
# Format: key_tag algorithm digest_type digest
```

### NAPTR Record
```hcl
data = "100 10 \"u\" \"E2U+sip\" \"!^.*$!sip:info@example.com!\" ."
priority = 100
# Format: order preference flags service regexp replacement
```

### URI Record
```hcl
data = "https://example.com/path"
priority = 10
weight = 1
```

### LOC Record
```hcl
data = "37 23 30.000 N 122 02 30.000 W 10m 100m 100m 10m"
# Format: lat lon alt size hp vp
```

## Validation Rules

### TTL Validation
- Minimum: 600 seconds (10 minutes)
- Maximum: 2,592,000 seconds (30 days)
- Default: 3600 seconds (1 hour)

### Priority Validation
- Range: 0-65535
- Lower values have higher priority

### Port Validation
- Range: 1-65535

### Weight Validation
- Range: 0-65535
- Used for load balancing

### Data Format Validation
- **IPv4**: Standard dotted decimal notation
- **IPv6**: Standard hexadecimal notation with colons
- **Domain names**: Valid FQDN format
- **Text**: Maximum 65535 characters

## Usage Examples

See the [examples/dns-record-types](examples/dns-record-types) directory for complete configuration examples of all supported record types.

## GoDaddy API Compatibility

All record types are compatible with the GoDaddy Domains API v1. The provider handles proper formatting and validation before sending requests to the API.