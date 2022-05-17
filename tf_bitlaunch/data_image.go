package tf_bitlaunch

import (
	"context"

	"github.com/bitlaunchio/gobitlaunch"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Holds details on Images and apps available when configuring a server. Matches https://developers.bitlaunch.io/reference/host-image-object",

		ReadContext: dataSourceImageRead,

		Schema: map[string]*schema.Schema{
			"host": {
				Description:  "Host Provider (DigitalOcean, Vultr, etc.)",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateHostID,
			},
			"distro_name": {
				Description:  "The name of the Linux Distibution or one-click app.",
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"distro_name", "version_name"},
			},
			"version_name": {
				Description:  "The Specific Image Version",
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"distro_name", "version_name"},
			},
			"type": {
				Description: "The type of the image: image or app.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"min_disk_size": {
				Description: "The minimum disk size available in GB.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"unavailable_regions": {
				Description: "A list of unavailable subregion IDs.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"extra_cost_per_month": {
				Description: "Extra monthly cost.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"is_windows": {
				Description: "Flag to determine if the image is Windows-based.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"password_unsupported": {
				Description: "If setting a password is supported.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func setDataImage(data *schema.ResourceData, image *gobitlaunch.HostImage, version *gobitlaunch.HostImageVersion, hostName string) diag.Diagnostics {
	var diags diag.Diagnostics

	data.SetId(version.ID)
	if err := data.Set("host", hostName); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("distro_name", image.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("version_name", version.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("type", image.Type); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("min_disk_size", image.MinDiskSize); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("unavailable_regions", image.UnavailableRegions); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("extra_cost_per_month", image.ExtraCostPerMonth); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("is_windows", image.Windows); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("password_unsupported", version.PasswordUnsupported); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func dataSourceImageRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*apiClient).client
	tflog.Trace(ctx, "Getting a Image")

	hostName := data.Get("host").(string)
	hostID := HostIDs[hostName]
	ops, err := client.CreateOptions.Show(hostID)
	if err != nil {
		return diag.FromErr(err)
	}

	name := data.Get("distro_name").(string)
	versionName := data.Get("version_name").(string)

	if len(name) == 0 && len(versionName) == 0 {
		return diag.Errorf("Require one of distro_name, version_name")
	}
	for _, image := range ops.Images {
		if len(name) != 0 && image.Name != name {
			continue
		}

		if len(versionName) == 0 || image.DefaultVersion.Description == versionName {
			// Get first/default version
			setDataImage(data, &image, &image.DefaultVersion, hostName)
			return diags
		}

		// Otherwise look in versions
		for _, version := range image.Versions {
			if version.Description == versionName {
				setDataImage(data, &image, &version, hostName)
				return diags
			}
		}
	}

	return diag.Errorf("Can't find matching Image")
}
