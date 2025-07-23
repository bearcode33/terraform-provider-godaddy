package godaddy

import (
	"fmt"
	"strconv"
	"strings"
)

// DNS record types supported by GoDaddy
const (
	RecordTypeA     = "A"
	RecordTypeAAAA  = "AAAA"
	RecordTypeCAA   = "CAA"
	RecordTypeCNAME = "CNAME"
	RecordTypeMX    = "MX"
	RecordTypeNS    = "NS"
	RecordTypePTR   = "PTR"
	RecordTypeSOA   = "SOA"
	RecordTypeSRV   = "SRV"
	RecordTypeTXT   = "TXT"
	RecordTypeDS    = "DS"
	RecordTypeSSHFP = "SSHFP"
	RecordTypeTLSA  = "TLSA"
	RecordTypeLOC   = "LOC"
	RecordTypeNAPTR = "NAPTR"
	RecordTypeURI   = "URI"
	RecordTypeCERT  = "CERT"
	RecordTypeDNAME = "DNAME"
)

// ValidDNSRecordTypes returns all valid DNS record types
func ValidDNSRecordTypes() []string {
	return []string{
		RecordTypeA,
		RecordTypeAAAA,
		RecordTypeCAA,
		RecordTypeCNAME,
		RecordTypeMX,
		RecordTypeNS,
		RecordTypePTR,
		RecordTypeSOA,
		RecordTypeSRV,
		RecordTypeTXT,
		RecordTypeDS,
		RecordTypeSSHFP,
		RecordTypeTLSA,
		RecordTypeLOC,
		RecordTypeNAPTR,
		RecordTypeURI,
		RecordTypeCERT,
		RecordTypeDNAME,
	}
}

// DNSRecordRequiredFields defines which fields are required for each record type
func DNSRecordRequiredFields(recordType string) map[string]bool {
	switch strings.ToUpper(recordType) {
	case RecordTypeA:
		return map[string]bool{"data": true}
	case RecordTypeAAAA:
		return map[string]bool{"data": true}
	case RecordTypeCAA:
		return map[string]bool{"data": true}
	case RecordTypeCNAME:
		return map[string]bool{"data": true}
	case RecordTypeMX:
		return map[string]bool{"data": true, "priority": true}
	case RecordTypeNS:
		return map[string]bool{"data": true}
	case RecordTypePTR:
		return map[string]bool{"data": true}
	case RecordTypeSOA:
		return map[string]bool{"data": true}
	case RecordTypeSRV:
		return map[string]bool{"data": true, "priority": true, "weight": true, "port": true}
	case RecordTypeTXT:
		return map[string]bool{"data": true}
	case RecordTypeDS:
		return map[string]bool{"data": true}
	case RecordTypeSSHFP:
		return map[string]bool{"data": true}
	case RecordTypeTLSA:
		return map[string]bool{"data": true}
	case RecordTypeLOC:
		return map[string]bool{"data": true}
	case RecordTypeNAPTR:
		return map[string]bool{"data": true, "priority": true}
	case RecordTypeURI:
		return map[string]bool{"data": true, "priority": true, "weight": true}
	case RecordTypeCERT:
		return map[string]bool{"data": true}
	case RecordTypeDNAME:
		return map[string]bool{"data": true}
	default:
		return map[string]bool{"data": true}
	}
}

// ValidateRecordData validates the data field based on record type
func ValidateRecordData(recordType, data string) error {
	switch strings.ToUpper(recordType) {
	case RecordTypeA:
		return validateIPv4(data)
	case RecordTypeAAAA:
		return validateIPv6(data)
	case RecordTypeCAA:
		return validateCAA(data)
	case RecordTypeCNAME:
		return validateDomainName(data)
	case RecordTypeMX:
		return validateDomainName(data)
	case RecordTypeNS:
		return validateDomainName(data)
	case RecordTypePTR:
		return validateDomainName(data)
	case RecordTypeSRV:
		return validateDomainName(data)
	case RecordTypeTXT:
		return validateTXT(data)
	case RecordTypeDS:
		return validateDS(data)
	case RecordTypeSSHFP:
		return validateSSHFP(data)
	case RecordTypeTLSA:
		return validateTLSA(data)
	case RecordTypeNAPTR:
		return validateNAPTR(data)
	case RecordTypeURI:
		return validateURI(data)
	}
	return nil
}

// Validation functions for different record types
func validateIPv4(ip string) error {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return fmt.Errorf("invalid IPv4 address: %s", ip)
	}
	for _, part := range parts {
		if num, err := strconv.Atoi(part); err != nil || num < 0 || num > 255 {
			return fmt.Errorf("invalid IPv4 address: %s", ip)
		}
	}
	return nil
}

func validateIPv6(ip string) error {
	// Basic IPv6 validation - could be more comprehensive
	if !strings.Contains(ip, ":") {
		return fmt.Errorf("invalid IPv6 address: %s", ip)
	}
	return nil
}

func validateCAA(data string) error {
	parts := strings.Fields(data)
	if len(parts) < 3 {
		return fmt.Errorf("CAA record must have format: flags tag value")
	}

	// Validate flags (0-255)
	if flags, err := strconv.Atoi(parts[0]); err != nil || flags < 0 || flags > 255 {
		return fmt.Errorf("CAA flags must be 0-255")
	}

	return nil
}

func validateDomainName(domain string) error {
	if len(domain) == 0 || len(domain) > 253 {
		return fmt.Errorf("invalid domain name length")
	}
	return nil
}

func validateTXT(data string) error {
	if len(data) > 65535 {
		return fmt.Errorf("TXT record too long (max 65535 characters)")
	}
	return nil
}

func validateDS(data string) error {
	parts := strings.Fields(data)
	if len(parts) != 4 {
		return fmt.Errorf("DS record must have format: key_tag algorithm digest_type digest")
	}
	return nil
}

func validateSSHFP(data string) error {
	parts := strings.Fields(data)
	if len(parts) != 3 {
		return fmt.Errorf("SSHFP record must have format: algorithm fp_type fingerprint")
	}
	return nil
}

func validateTLSA(data string) error {
	parts := strings.Fields(data)
	if len(parts) != 4 {
		return fmt.Errorf("TLSA record must have format: cert_usage selector matching_type cert_data")
	}
	return nil
}

func validateNAPTR(data string) error {
	parts := strings.Fields(data)
	if len(parts) < 6 {
		return fmt.Errorf("NAPTR record must have format: order preference flags service regexp replacement")
	}
	return nil
}

func validateURI(data string) error {
	if !strings.HasPrefix(data, "http://") && !strings.HasPrefix(data, "https://") {
		return fmt.Errorf("URI record must start with http:// or https://")
	}
	return nil
}
