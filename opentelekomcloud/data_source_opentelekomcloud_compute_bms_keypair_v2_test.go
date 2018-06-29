package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBMSV2KeypairDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccBmsKeyPairPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccBMSV2KeypairDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBMSV2KeypairDataSourceID("data.opentelekomcloud_compute_bms_keypairs_v2.keypair"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_compute_bms_keypairs_v2.keypair", "name", OS_KEYPAIR_NAME),
				),
			},
		},
	})
}

func testAccCheckBMSV2KeypairDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find keypair data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Keypair data source ID not set")
		}

		return nil
	}
}

var testAccBMSV2KeypairDataSource_basic = fmt.Sprintf(`
data "opentelekomcloud_compute_bms_keypairs_v2" "keypair" {
  name = "%s"
}
`, OS_KEYPAIR_NAME)
