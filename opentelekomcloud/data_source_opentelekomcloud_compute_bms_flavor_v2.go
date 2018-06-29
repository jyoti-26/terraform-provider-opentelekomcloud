package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/bms/v2/flavors"
	"log"
	"strings"
)

func dataSourceBMSFlavorV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBMSFlavorV2Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"min_ram": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"ram": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"min_disk": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"disk": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"swap": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"rx_tx_factor": {
				Type:     schema.TypeFloat,
				Optional: true,
			},

			// Computed values
			"is_public": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "id",
			},

			"sort_dir": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "asc",
				ValidateFunc: dataSourceImagesImageV2SortDirection,
			},
		},
	}
}

func dataSourceBMSFlavorV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	flavorClient, err := config.bmsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Opentelekom bms client: %s", err)
	}

	listOpts := flavors.ListOpts{
		MinDisk:    d.Get("min_disk").(int),
		MinRAM:     d.Get("min_ram").(int),
		AccessType: flavors.PublicAccess,
		Name:       d.Get("name").(string),
		ID:         d.Get("id").(string),
		SortKey:    d.Get("sort_key").(string),
		SortDir:    d.Get("sort_dir").(string),
	}
	var flavor flavors.Flavor
	refinedflavors, err := flavors.List(flavorClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve flavors: %s", err)
	}

	if len(refinedflavors) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedflavors) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	} else {
		flavor = refinedflavors[0]
	}

	//flavor := refinedflavors[0]

	if !strings.Contains(flavor.ID, "physical") {
		return fmt.Errorf("Flavors name starting with 'physical' are BMS flavors not: %s ", flavor.ID)
	}

	log.Printf("[DEBUG] Single Flavor found: %s", flavor.ID)
	d.SetId(flavor.ID)
	d.Set("name", flavor.Name)
	d.Set("disk", flavor.Disk)
	d.Set("min_disk", flavor.MinDisk)
	d.Set("sort_key", flavor.SortKey)
	d.Set("sort_dir", flavor.SortDir)
	d.Set("min_ram", flavor.MinRAM)
	d.Set("ram", flavor.RAM)
	d.Set("rx_tx_factor", flavor.RxTxFactor)
	d.Set("swap", flavor.Swap)
	d.Set("vcpus", flavor.VCPUs)
	d.Set("is_public", flavor.IsPublic)

	return nil
}
