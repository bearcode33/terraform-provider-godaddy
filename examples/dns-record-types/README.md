# DNS Record Types Examples

This directory contains examples for all DNS record types supported by the GoDaddy Terraform provider.

## Supported Record Types

### Basic Records
- **A**: IPv4 address mapping
- **AAAA**: IPv6 address mapping  
- **CNAME**: Canonical name (alias)
- **TXT**: Text records (SPF, DKIM, DMARC, etc.)

### Mail Records
- **MX**: Mail exchange servers

### Service Records
- **SRV**: Service location records
- **NS**: Name server delegation
- **PTR**: Reverse DNS lookups

### Security Records
- **CAA**: Certificate Authority Authorization
- **SSHFP**: SSH key fingerprints
- **TLSA**: TLS certificate association
- **DS**: DNSSEC Delegation Signer

### Advanced Records
- **NAPTR**: Naming Authority Pointer
- **URI**: Uniform Resource Identifier
- **LOC**: Geographic location
- **CERT**: Certificate records
- **DNAME**: Delegation Name

## Usage

1. Copy the desired record configuration from `all-record-types.tf`
2. Modify the domain name and record values for your use case
3. Apply with Terraform:

```bash
terraform init
terraform plan
terraform apply
```

## Record-Specific Requirements

### MX Records
- Require `priority` field (0-65535)
- Lower priority values are preferred

### SRV Records
- Require `priority`, `weight`, `port`, `service`, and `protocol` fields
- Service and protocol must start with underscore (e.g., `_sip`, `_tcp`)

### CAA Records
- Format: `flags tag value`
- Common tags: `issue`, `issuewild`, `iodef`

### SSHFP Records
- Format: `algorithm fingerprint_type fingerprint`
- Algorithm: 1 (RSA), 2 (DSA), 3 (ECDSA), 4 (Ed25519)
- Fingerprint type: 1 (SHA-1), 2 (SHA-256)

### TLSA Records
- Format: `cert_usage selector matching_type cert_data`
- Used for DANE (DNS-based Authentication of Named Entities)

### DS Records
- Format: `key_tag algorithm digest_type digest`
- Used for DNSSEC delegation

## Validation

The provider includes built-in validation for:
- Record type validity
- Required fields based on record type
- Data format validation (IPv4, IPv6, domain names, etc.)
- Value ranges (TTL, priority, port, weight)

## TTL Considerations

- Minimum TTL: 600 seconds (10 minutes)
- Maximum TTL: 2592000 seconds (30 days)
- Default TTL: 3600 seconds (1 hour)
- Lower TTL values provide faster updates but higher DNS query load