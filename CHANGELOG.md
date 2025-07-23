# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-07-23

### Added
- Initial release of the GoDaddy Terraform Provider
- Complete GoDaddy API client implementation with authentication
- Domain resource for managing domain settings (locked, privacy, auto-renewal, nameservers)
- DNS record resource supporting all 18+ DNS record types:
  - A, AAAA, CAA, CNAME, MX, NS, PTR, SOA, SRV, TXT
  - DS, SSHFP, TLSA, LOC, NAPTR, URI, CERT, DNAME
- Data sources for querying domains and DNS records
- Comprehensive validation for all DNS record types
- Custom validators for TTL, priority, weight, and port values
- Support for both production and test (OTE) environments
- Complete documentation and examples
- Terraform Plugin Framework v1.15 implementation
- Import/export functionality for existing resources

### Features
- **Domain Management**: Lock/unlock domains, configure privacy settings, manage auto-renewal
- **DNS Management**: Full support for all DNS record types with proper validation
- **Environment Support**: Switch between production and test environments
- **Authentication**: API key/secret authentication with environment variable support
- **Validation**: Type-specific validation for DNS records (e.g., IP addresses, domains, emails)
- **Error Handling**: Comprehensive error handling with helpful messages

### Documentation
- Complete provider documentation with authentication setup
- Resource and data source documentation with all attributes
- Extensive examples for all DNS record types
- Development guide and contribution instructions
- Terraform Registry manifest for publication

### Examples
- Basic domain and DNS record management
- Complete setup with multiple record types
- Advanced configurations with custom nameservers
- Email security setup (SPF, DKIM, DMARC)
- SSL certificate configuration (CAA records)
