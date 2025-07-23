package godaddy

import (
	"testing"
)

func TestValidDNSRecordTypes(t *testing.T) {
	validTypes := ValidDNSRecordTypes()

	expectedTypes := []string{
		"A", "AAAA", "CAA", "CNAME", "MX", "NS", "PTR", "SOA", "SRV", "TXT",
		"DS", "SSHFP", "TLSA", "LOC", "NAPTR", "URI", "CERT", "DNAME",
	}

	if len(validTypes) != len(expectedTypes) {
		t.Errorf("Expected %d valid types, got %d", len(expectedTypes), len(validTypes))
	}

	// Check that all expected types are present
	typeMap := make(map[string]bool)
	for _, vt := range validTypes {
		typeMap[vt] = true
	}

	for _, expected := range expectedTypes {
		if !typeMap[expected] {
			t.Errorf("Expected type %s not found in valid types", expected)
		}
	}
}

func TestDNSRecordRequiredFields(t *testing.T) {
	tests := []struct {
		recordType     string
		expectedFields map[string]bool
	}{
		{
			recordType:     "A",
			expectedFields: map[string]bool{"data": true},
		},
		{
			recordType:     "MX",
			expectedFields: map[string]bool{"data": true, "priority": true},
		},
		{
			recordType:     "SRV",
			expectedFields: map[string]bool{"data": true, "priority": true, "weight": true, "port": true},
		},
		{
			recordType:     "TXT",
			expectedFields: map[string]bool{"data": true},
		},
	}

	for _, test := range tests {
		t.Run(test.recordType, func(t *testing.T) {
			fields := DNSRecordRequiredFields(test.recordType)

			for field, required := range test.expectedFields {
				if fields[field] != required {
					t.Errorf("For %s record, expected field %s to be %v, got %v",
						test.recordType, field, required, fields[field])
				}
			}
		})
	}
}

func TestValidateRecordData(t *testing.T) {
	tests := []struct {
		name       string
		recordType string
		data       string
		expectErr  bool
	}{
		{
			name:       "Valid A record",
			recordType: "A",
			data:       "192.0.2.1",
			expectErr:  false,
		},
		{
			name:       "Invalid A record",
			recordType: "A",
			data:       "256.0.0.1",
			expectErr:  true,
		},
		{
			name:       "Valid IPv6",
			recordType: "AAAA",
			data:       "2001:db8::1",
			expectErr:  false,
		},
		{
			name:       "Invalid IPv6",
			recordType: "AAAA",
			data:       "192.0.2.1",
			expectErr:  true,
		},
		{
			name:       "Valid CNAME",
			recordType: "CNAME",
			data:       "example.com",
			expectErr:  false,
		},
		{
			name:       "Valid TXT",
			recordType: "TXT",
			data:       "v=spf1 include:_spf.google.com ~all",
			expectErr:  false,
		},
		{
			name:       "Valid CAA",
			recordType: "CAA",
			data:       "0 issue \"letsencrypt.org\"",
			expectErr:  false,
		},
		{
			name:       "Invalid CAA",
			recordType: "CAA",
			data:       "invalid caa",
			expectErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateRecordData(test.recordType, test.data)

			if test.expectErr && err == nil {
				t.Errorf("Expected error for %s record with data %s, but got none", test.recordType, test.data)
			}

			if !test.expectErr && err != nil {
				t.Errorf("Expected no error for %s record with data %s, but got: %v", test.recordType, test.data, err)
			}
		})
	}
}

func TestValidateIPv4(t *testing.T) {
	tests := []struct {
		ip        string
		expectErr bool
	}{
		{"192.0.2.1", false},
		{"0.0.0.0", false},
		{"255.255.255.255", false},
		{"256.0.0.1", true},
		{"192.0.2", true},
		{"192.0.2.1.1", true},
		{"not.an.ip", true},
	}

	for _, test := range tests {
		t.Run(test.ip, func(t *testing.T) {
			err := validateIPv4(test.ip)

			if test.expectErr && err == nil {
				t.Errorf("Expected error for IP %s, but got none", test.ip)
			}

			if !test.expectErr && err != nil {
				t.Errorf("Expected no error for IP %s, but got: %v", test.ip, err)
			}
		})
	}
}

func TestValidateCAA(t *testing.T) {
	tests := []struct {
		data      string
		expectErr bool
	}{
		{"0 issue \"letsencrypt.org\"", false},
		{"0 issuewild \"letsencrypt.org\"", false},
		{"128 iodef \"mailto:admin@example.com\"", false},
		{"300 issue \"letsencrypt.org\"", true}, // Invalid flags
		{"0 issue", true},                       // Missing value
		{"invalid", true},                       // Invalid format
	}

	for _, test := range tests {
		t.Run(test.data, func(t *testing.T) {
			err := validateCAA(test.data)

			if test.expectErr && err == nil {
				t.Errorf("Expected error for CAA data %s, but got none", test.data)
			}

			if !test.expectErr && err != nil {
				t.Errorf("Expected no error for CAA data %s, but got: %v", test.data, err)
			}
		})
	}
}
