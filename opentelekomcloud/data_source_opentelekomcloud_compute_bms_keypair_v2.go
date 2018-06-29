package opentelekomcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/bms/v2/keypairs"
	"log"
)

func dataSourceBMSKeypairV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBMSKeypairV2Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBMSKeypairV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	bmsClient, err := config.bmsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Opentelekom bms client: %s", err)
	}

	listOpts := keypairs.ListOpts{
		Name: d.Get("name").(string),
	}

	refinedKeypairs, err := keypairs.List(bmsClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve keypairs: %s", err)
	}

	if len(refinedKeypairs) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedKeypairs) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Vpc := refinedKeypairs[0]

	log.Printf("[INFO] Retrieved Keypairs using given filter %s: %+v", Vpc.Name, Vpc)
	d.SetId(Vpc.Name)

	d.Set("name", Vpc.Name)
	d.Set("public_key", Vpc.PublicKey)
	d.Set("fingerprint", Vpc.Fingerprint)
	d.Set("region", GetRegion(d, config))

	return nil
}
