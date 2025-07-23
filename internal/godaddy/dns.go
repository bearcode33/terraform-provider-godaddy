package godaddy

import (
	"context"
	"fmt"
)

func (c *Client) GetDNSRecords(ctx context.Context, domain string) ([]DNSRecord, error) {
	var result []DNSRecord
	err := c.Get(ctx, fmt.Sprintf("/v1/domains/%s/records", domain), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get DNS records for domain %s: %w", domain, err)
	}
	return result, nil
}

func (c *Client) GetDNSRecordsByType(ctx context.Context, domain, recordType string) ([]DNSRecord, error) {
	var result []DNSRecord
	err := c.Get(ctx, fmt.Sprintf("/v1/domains/%s/records/%s", domain, recordType), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get DNS records of type %s for domain %s: %w", recordType, domain, err)
	}
	return result, nil
}

func (c *Client) GetDNSRecordsByTypeAndName(ctx context.Context, domain, recordType, name string) ([]DNSRecord, error) {
	var result []DNSRecord
	err := c.Get(ctx, fmt.Sprintf("/v1/domains/%s/records/%s/%s", domain, recordType, name), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get DNS record %s.%s for domain %s: %w", name, recordType, domain, err)
	}
	return result, nil
}

func (c *Client) ReplaceDNSRecords(ctx context.Context, domain string, records []DNSRecord) error {
	err := c.Put(ctx, fmt.Sprintf("/v1/domains/%s/records", domain), records)
	if err != nil {
		return fmt.Errorf("failed to replace DNS records for domain %s: %w", domain, err)
	}
	return nil
}

func (c *Client) ReplaceDNSRecordsByType(ctx context.Context, domain, recordType string, records []DNSRecord) error {
	err := c.Put(ctx, fmt.Sprintf("/v1/domains/%s/records/%s", domain, recordType), records)
	if err != nil {
		return fmt.Errorf("failed to replace DNS records of type %s for domain %s: %w", recordType, domain, err)
	}
	return nil
}

func (c *Client) ReplaceDNSRecordsByTypeAndName(ctx context.Context, domain, recordType, name string, records []DNSRecord) error {
	err := c.Put(ctx, fmt.Sprintf("/v1/domains/%s/records/%s/%s", domain, recordType, name), records)
	if err != nil {
		return fmt.Errorf("failed to replace DNS record %s.%s for domain %s: %w", name, recordType, domain, err)
	}
	return nil
}

func (c *Client) AddDNSRecord(ctx context.Context, domain string, record DNSRecord) error {
	err := c.Patch(ctx, fmt.Sprintf("/v1/domains/%s/records", domain), []DNSRecord{record})
	if err != nil {
		return fmt.Errorf("failed to add DNS record to domain %s: %w", domain, err)
	}
	return nil
}

func (c *Client) DeleteDNSRecord(ctx context.Context, domain, recordType, name string) error {
	err := c.Delete(ctx, fmt.Sprintf("/v1/domains/%s/records/%s/%s", domain, recordType, name))
	if err != nil {
		return fmt.Errorf("failed to delete DNS record %s.%s from domain %s: %w", name, recordType, domain, err)
	}
	return nil
}
