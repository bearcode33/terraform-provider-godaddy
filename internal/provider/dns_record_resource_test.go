package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDNSRecordResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDNSRecordResourceConfig("test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("godaddy_dns_record.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("godaddy_dns_record.test", "type", "A"),
					resource.TestCheckResourceAttr("godaddy_dns_record.test", "name", "test"),
					resource.TestCheckResourceAttr("godaddy_dns_record.test", "data", "192.0.2.1"),
					resource.TestCheckResourceAttr("godaddy_dns_record.test", "ttl", "3600"),
					resource.TestCheckResourceAttrSet("godaddy_dns_record.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "godaddy_dns_record.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return "example.com/A/test/192.0.2.1", nil
				},
			},
			// Update and Read testing
			{
				Config: testAccDNSRecordResourceConfig("test-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("godaddy_dns_record.test", "data", "192.0.2.2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDNSRecordResourceConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "godaddy_dns_record" "test" {
  domain = "example.com"
  type   = "A"
  name   = "%s"
  data   = "192.0.2.%s"
  ttl    = 3600
}
`, providerConfig, name, func() string {
		if name == "test-updated" {
			return "2"
		}
		return "1"
	}())
}
