package tf_bitlaunch

import (
	"context"
	"strings"

	"github.com/bitlaunchio/gobitlaunch"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRegion() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Holds available region configurations for a server. Matches https://developers.bitlaunch.io/reference/host-region-object",

		ReadContext: dataSourceRegionRead,

		Schema: map[string]*schema.Schema{
			"host": {
				Description:  "Host Provider (DigitalOcean, Vultr, etc.)",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateHostID,
			},
			"region_name": {
				Description:  "The name of the Region.",
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"region_name", "slug"},
			},
			"slug": {
				Description:  "The Specific Subregion slug.",
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"region_name", "slug"},
			},
			"iso": {
				Description: "The ISO code for the region.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"unavailable_sizes": {
				Description: "A list of the unavailable sizes for this subregion.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
		},
	}
}

func setDataRegion(data *schema.ResourceData, region *gobitlaunch.HostRegion, subregion *gobitlaunch.HostSubRegion, hostName string) diag.Diagnostics {
	var diags diag.Diagnostics

	data.SetId(subregion.ID)
	if err := data.Set("host", hostName); err != nil {
		return diag.FromErr(err)
	}
	// The API doesn't trim some of these...
	if err := data.Set("region_name", strings.TrimSpace(region.Name)); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("slug", subregion.Slug); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("iso", region.ISO); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("unavailable_sizes", subregion.UnavailableSizes); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func dataSourceRegionRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*apiClient).client
	tflog.Trace(ctx, "Getting a Region")

	hostName := data.Get("host").(string)
	hostID := HostIDs[hostName]
	ops, err := client.CreateOptions.Show(hostID)
	if err != nil {
		return diag.FromErr(err)
	}

	// The API doesn't trim some of these...
	name := strings.TrimSpace(data.Get("region_name").(string))
	slug := data.Get("slug").(string)

	if len(name) == 0 && len(slug) == 0 {
		return diag.Errorf("Require one of region_name, slug")
	}
	for _, region := range ops.Regions {
		if len(name) != 0 && strings.TrimSpace(region.Name) != name {
			continue
		}

		if len(slug) == 0 || region.DefaultSubregion.Slug == slug {
			// Get first/default subregion
			setDataRegion(data, &region, &region.DefaultSubregion, hostName)
			return diags
		}

		// Otherwise look in subregions
		for _, subregion := range region.Subregions {
			if subregion.Slug == slug {
				setDataRegion(data, &region, &subregion, hostName)
				return diags
			}
		}
	}

	return diag.Errorf("Can't find matching Region")
}
