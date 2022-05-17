package tf_bitlaunch

import (
	"context"

	"github.com/bitlaunchio/gobitlaunch"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSize() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Holds details on available size configurations for a server. Matches https://developers.bitlaunch.io/reference/host-size-object",

		ReadContext: dataSourceSizeRead,

		Schema: map[string]*schema.Schema{
			"host": {
				Description:  "Host Provider (DigitalOcean, Vultr, etc.)",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateHostID,
			},
			"cpu_count": {
				Description:  "The amount of vCPU's included.",
				Type:         schema.TypeInt,
				Optional:     true,
				AtLeastOneOf: []string{"cpu_count", "disk_gb", "memory_mb"},
			},
			"disk_gb": {
				Description:  "The amount of disk space included.",
				Type:         schema.TypeInt,
				Optional:     true,
				AtLeastOneOf: []string{"cpu_count", "disk_gb", "memory_mb"},
			},
			"memory_mb": {
				Description:  "The amount of memory (RAM) included.",
				Type:         schema.TypeInt,
				Optional:     true,
				AtLeastOneOf: []string{"cpu_count", "disk_gb", "memory_mb"},
			},
			"slug": {
				Description: "A human readable string.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"bandwidth_gb": {
				Description: "The available monthly bandwidth in GB.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"cost_per_hour": {
				Description: "The amount of balance deducted per hour.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"cost_per_month": {
				Description: "The amount in USD charged per month.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"plan_type": {
				Description: "Some hosts offer a different plan type for different usage. You should refer to the host documentation for more information.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"disks": {
				Description: "Details on disks included with the size.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "The type of storage disk (SSD/HDD).",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"count": {
							Description: "The amount of disks.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"size": {
							Description: "The size of the disk(s).",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"unit": {
							Description: "The unit of measurement for the disk size.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func setDataSize(data *schema.ResourceData, size *gobitlaunch.HostSize, hostName string) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := data.Set("host", hostName); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("cpu_count", size.CPUCount); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("disk_gb", size.DiskGB); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("memory_mb", size.MemoryMB); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("slug", size.Slug); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("bandwidth_gb", size.BandwidthGB); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("cost_per_hour", size.CostPerHour); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("cost_per_month", size.CostPerMonth); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("plan_type", size.PlanType); err != nil {
		return diag.FromErr(err)
	}

	disks := make([]interface{}, len(size.Disks), len(size.Disks))
	for i, disk := range size.Disks {
		tfDisk := make(map[string]interface{})
		tfDisk["type"] = disk.Type
		tfDisk["count"] = disk.Count
		tfDisk["size"] = disk.Size
		tfDisk["unit"] = disk.Unit
		disks[i] = tfDisk
	}
	if err := data.Set("disks", disks); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func dataSourceSizeRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*apiClient).client
	tflog.Trace(ctx, "Getting a Size")

	hostName := data.Get("host").(string)
	hostID := HostIDs[hostName]
	ops, err := client.CreateOptions.Show(hostID)
	if err != nil {
		return diag.FromErr(err)
	}

	cpuCount := data.Get("cpu_count").(int)
	diskGB := data.Get("disk_gb").(int)
	memoryMB := data.Get("memory_mb").(int)
	if cpuCount == 0 && diskGB == 0 && memoryMB == 0 {
		return diag.Errorf("Require one of cpu_count, disk_gb, memory_mb")
	}
	for _, size := range ops.Sizes {
		if cpuCount != 0 && size.CPUCount != cpuCount {
			continue
		}
		if diskGB != 0 && size.DiskGB != diskGB {
			continue
		}
		if memoryMB != 0 && size.MemoryMB != memoryMB {
			continue
		}
		// Found the first matching?
		data.SetId(size.ID)
		setDataSize(data, &size, hostName)
		return diags
	}

	return diag.Errorf("Can't find matching Size")
}
