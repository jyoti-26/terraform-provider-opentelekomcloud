package servers

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
	"reflect"
)

// ListServerOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned. Marker and Limit are used
// for pagination.
type ListServerOpts struct {
	// ID uniquely identifies this server amongst all other servers,
	// including those not accessible to the current tenant.
	ID string `json:"id"`
	//ID of the user to which the BMS belongs.
	UserID string `json:"user_id"`
	//Contains the nova-compute status
	HostStatus string `json:"host_status"`
	//Contains the host ID of the BMS.
	HostID string `json:"hostid"`
	// KeyName indicates which public key was injected into the server on launch.
	KeyName string `json:"key_name"`
	// Specifies the BMS name.
	Name    string `q:"name"`
	ImageID string `q:"image"`
	// Specifies flavor ID.
	FlavorID string `q:"flavor"`
	// Specifies the BMS status.
	Status string `q:"status"`
	//Filters out the BMSs that have been updated since the changes-since time.
	// The parameter is in ISO 8601 time format, for example, 2013-06-09T06:42:18Z.
	ChangesSince string `q:"changes-since"`
	//Specifies whether to query the BMSs of all tenants. This parameter is available only to administrators.
	// The value can be 0 (do not query the BMSs of all tenants) or 1 (query the BMSs of all tenants).
	AllTenants int `q:"all_tenants"`
	//Specifies the IP address. This parameter supports fuzzy matching.
	IP string `q:"ip"`
	//Specifies the tag list. Returns BMSs that match all tags. Use commas (,) to separate multiple tags
	Tags string `q:"tags"`
	//Specifies the tag list. Returns BMSs that match any tag
	TagsAny string `q:"tags-any"`
	//Specifies the tag list. Returns BMSs that do not match all tags.
	NotTags string `q:"not-tags"`
	//Specifies the tag list. Returns BMSs that do not match any of the tags.
	NotTagsAny int `q:"not-tags-any"`
	//Specifies the BMS sorting attribute, which can be the BMS UUID (uuid), BMS status (vm_state),
	// BMS name (display_name), BMS task status (task_state), power status (power_state),
	// creation time (created_at), last time when the BMS is updated (updated_at), and availability zone
	// (availability_zone). You can specify multiple sort_key and sort_dir pairs.
	SortKey string `q:"sort_key"`
	//Specifies the sorting direction, i.e. asc or desc.
	SortDir string `q:"sort_dir"`
}

func FilterParam(opts ListServerOpts) (filter ListServerOpts) {

	if opts.ID != "" {
		filter.ID = opts.ID
	}
	filter.Name = opts.Name
	filter.Status = opts.Status
	filter.FlavorID = opts.FlavorID
	filter.ChangesSince = opts.ChangesSince
	filter.SortKey = opts.SortKey
	filter.SortDir = opts.SortDir
	filter.AllTenants = opts.AllTenants
	filter.IP = opts.IP
	filter.Tags = opts.Tags
	filter.TagsAny = opts.TagsAny
	filter.NotTags = opts.NotTags
	filter.NotTagsAny = opts.NotTagsAny
	filter.ImageID = opts.ImageID

	return filter
}

// ListServer returns a Pager which allows you to iterate over a collection of
// dedicated hosts Server resources. It accepts a ListServerOpts struct, which allows you to
// filter the returned collection for greater efficiency.
func ListServer(c *golangsdk.ServiceClient, opts ListServerOpts) ([]Server, error) {
	c.Microversion = "2.26"
	filter := FilterParam(opts)
	q, err := golangsdk.BuildQueryString(&filter)
	if err != nil {
		return nil, err
	}
	u := listDetailURL(c) + q.String()
	pages, err := pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return ServerPage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()

	allservers, err := ExtractServers(pages)
	if err != nil {
		return nil, err
	}
	return FilterServers(allservers, opts)
}

func FilterServers(servers []Server, opts ListServerOpts) ([]Server, error) {
	var refinedServers []Server
	var matched bool
	m := map[string]interface{}{}

	if opts.ID != "" {
		m["ID"] = opts.ID
	}
	if opts.UserID != "" {
		m["UserID"] = opts.UserID
	}
	if opts.HostStatus != "" {
		m["HostStatus"] = opts.HostStatus
	}
	if opts.HostID != "" {
		m["HostID"] = opts.HostID
	}
	if opts.KeyName != "" {
		m["KeyName"] = opts.KeyName
	}
	if len(m) > 0 && len(servers) > 0 {
		for _, server := range servers {
			matched = true

			for key, value := range m {
				if sVal := getStructServerField(&server, key); !(sVal == value) {
					matched = false
				}
			}
			if matched {
				refinedServers = append(refinedServers, server)
			}
		}
	} else {
		refinedServers = servers
	}

	return refinedServers, nil
}

func getStructServerField(v *Server, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

// Get requests details on a single server, by ID.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: map[string]string{"X-OpenStack-Nova-API-Version": "2.26"},
	})
	return
}
