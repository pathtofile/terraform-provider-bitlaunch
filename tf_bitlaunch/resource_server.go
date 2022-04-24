package tf_bitlaunch

import (
	"context"
	"fmt"
	"time"

	"github.com/bitlaunchio/gobitlaunch"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// https://developers.bitlaunch.io/reference/create-server
func resourceServer() *schema.Resource {
	return &schema.Resource{
		Description: "Virtual Machine Server",

		CreateContext: resourceServerCreate,
		ReadContext:   resourceServerRead,
		DeleteContext: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"host": {
				Description: "The host for the server to reside on.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The name of the server.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"image_id": {
				Description: "The image ID to use on the server.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"size_id": {
				Description: "The size ID of the server to be provisioned to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"region_id": {
				Description: "The region ID of the location that the server will reside at.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"ssh_keys": {
				Description: "An array of SSH key IDs to place on the server for authentication. Must be used if no password is designated of if the selected image does not support passwords.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				ForceNew:    true,
			},
			"password": {
				Description: "The root user password to set on the server. Must be used if no SSH keys designated.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"initscript": {
				Description: "A script to run on first boot of the server. Only hosts with initScript enabled can use this feature.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"wait_for_ip": {
				Description: "Wait to get IP Address",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			"ipv4": {
				Description: "The name of the key.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status": {
				Description: "The name of the key.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created": {
				Description: "The creation date of the server.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"image_description": {
				Description: "The description of the image installed on the server.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"rate": {
				Description: "The hourly rate of the server that will be deducted from your account balance every hour.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func setDataServer(data *schema.ResourceData, server *gobitlaunch.Server, hostName string) diag.Diagnostics {
	var diags diag.Diagnostics

	data.SetId(server.ID)
	if err := data.Set("host", hostName); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("name", server.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("image_id", server.Image); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("image_description", server.ImageDesc); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("size_id", server.Size); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("region_id", server.Region); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("ipv4", server.Ipv4); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("status", server.Status); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("created", server.Created.Format(time.RFC3339)); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("rate", server.Rate); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceServerCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*apiClient).client
	tflog.Trace(ctx, "Creating an server")

	hostName := data.Get("host").(string)
	hostID := HostIDs[hostName]

	server := gobitlaunch.CreateServerOptions{
		HostID:      hostID,
		Name:        data.Get("name").(string),
		HostImageID: data.Get("image_id").(string),
		SizeID:      data.Get("size_id").(string),
		RegionID:    data.Get("region_id").(string),
	}

	// Get SSH Keys from list
	sshKeysRaw := data.Get("ssh_keys").([]interface{})
	sshKeys := make([]string, len(sshKeysRaw))
	for i, raw := range sshKeysRaw {
		sshKeys[i] = raw.(string)
	}
	if len(sshKeys) > 0 {
		server.SSHKeys = sshKeys
	}

	password := data.Get("password").(string)
	if len(password) > 0 {
		server.Password = password
	}
	initScript := data.Get("initscript").(string)
	if len(initScript) > 0 {
		server.InitScript = initScript
	}

	newServer, err := client.Server.Create(&server)
	if err != nil {
		return diag.FromErr(err)
	}

	if data.Get("wait_for_ip").(bool) {
		// Poll until server.Status == "ok"
		maxTime := time.Now().Add(60 * time.Second)
		for {
			if time.Now().After(maxTime) {
				return diag.Errorf("Timed out getting IPv4 address")
			}

			newServer, err = client.Server.Show(newServer.ID)
			if err != nil {
				return diag.FromErr(err)
			}
			if newServer.Status == "ok" {
				// Server IPv4 should now be there?
				break
			}
			if newServer.Status == "error" || newServer.Status == "stopped" {
				return diag.Errorf("Server creation returned error status")
			}
			time.Sleep(1 * time.Second)
		}
	}

	setDataServer(data, newServer, hostName)
	tflog.Trace(ctx, fmt.Sprintf("created Server %s", newServer.ID))

	return diags
}

func resourceServerRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*apiClient).client
	hostName := data.Get("host").(string)

	tflog.Trace(ctx, "Reading a server")

	servers, err := client.Server.List()
	if err != nil {
		return diag.FromErr(err)
	}
	for _, server := range servers {
		if server.ID == data.Id() {
			setDataServer(data, &server, hostName)
			return diags
		}
	}

	tflog.Trace(ctx, "server not found")
	data.SetId("")
	return diags
}

func resourceServerDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*apiClient).client
	tflog.Trace(ctx, "Deleting a server")

	err := client.Server.Destroy(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
