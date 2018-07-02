package opentelekomcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/bms/v2/servers"
	"log"
)

func dataSourceBMSServersV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBMSServersV2Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"progress": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"key_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_links": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rel": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"href": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_links": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rel": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"href": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"access_ip_v4": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_ip_v6": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rel": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"href": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"admin_pass": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"addresses": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": &schema.Schema{
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
			"all_tenants": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_groups": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tags": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"locked": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"config_drive": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"kernel_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"hypervisor_hostname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBMSServersV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	bmsClient, err := config.bmsClient(GetRegion(d, config))

	listServerOpts := servers.ListServerOpts{
		ID:       d.Get("id").(string),
		Name:     d.Get("name").(string),
		Status:   d.Get("status").(string),
		UserID:   d.Get("user_id").(string),
		KeyName:  d.Get("key_name").(string),
		FlavorID: d.Get("flavor_id").(string),
		ImageID:  d.Get("image_id").(string),
		Tags:     d.Get("tag").(string),
	}
	pages, err := servers.ListServer(bmsClient, listServerOpts)

	if err != nil {
		return fmt.Errorf("Unable to retrieve deh server: %s", err)
	}

	if len(pages) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	if len(pages) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}
	server := pages[0]

	log.Printf("[INFO] Retrieved Deh Server using given filter %s: %+v", server.ID, server)
	d.SetId(server.ID)

	var image []map[string]interface{}
	for _, value := range server.Image.Links {
		mapping := map[string]interface{}{
			"rel":  value.Rel,
			"href": value.Href,
		}
		image = append(image, mapping)
	}

	var flavor []map[string]interface{}
	for _, value := range server.Flavor.Links {
		mapping := map[string]interface{}{
			"rel":  value.Rel,
			"href": value.Href,
		}
		flavor = append(flavor, mapping)
	}

	var link []map[string]interface{}
	for _, value := range server.Links {
		mapping := map[string]interface{}{
			"rel":  value.Rel,
			"href": value.Href,
		}
		link = append(link, mapping)
	}

	var secGroups []map[string]interface{}
	for _, value := range server.SecurityGroups {
		mapping := map[string]interface{}{
			"name": value.Name,
		}
		secGroups = append(secGroups, mapping)
	}

	d.Set("server_id", server.ID)
	d.Set("user_id", server.UserID)
	d.Set("name", server.Name)
	d.Set("status", server.Status)
	d.Set("flavor_id", server.Flavor.ID)
	d.Set("flavor_links", flavor)
	d.Set("addresses", server.Addresses)
	d.Set("metadata", server.Metadata)
	d.Set("tenant_id", server.TenantID)
	d.Set("image_id", server.Image.ID)
	d.Set("image_links", image)
	d.Set("access_ip_v4", server.AccessIPv4)
	d.Set("access_ip_v6", server.AccessIPv6)
	d.Set("access_ip_v6", server.KeyName)
	d.Set("links", link)
	d.Set("admin_pass", server.AdminPass)
	d.Set("progress", server.Progress)
	d.Set("key_name", server.KeyName)
	d.Set("security_groups", secGroups)
	d.Set("locked", server.Locked)
	d.Set("config_drive", server.ConfigDrive)
	d.Set("availability_zone", server.AvailabilityZone)
	d.Set("description", server.Description)
	d.Set("kernel_id", server.KernelId)
	d.Set("hypervisor_hostname", server.HypervisorHostName)
	d.Set("instance_name", server.InstanceName)
	d.Set("tags", server.Tags)
	d.Set("region", GetRegion(d, config))
	networks, err := flattenServerNetwork(d, meta)
	if err != nil {
		return err
	}
	if err := d.Set("addresses", networks); err != nil {
		return fmt.Errorf("[DEBUG] Error saving network to state for OpenTelekomCloud server (%s): %s", d.Id(), err)
	}

	return nil
}
