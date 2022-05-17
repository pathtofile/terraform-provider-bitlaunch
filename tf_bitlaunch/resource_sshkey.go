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

// https://developers.bitlaunch.io/reference/create-ssh-key
func resourceSSHKey() *schema.Resource {
	return &schema.Resource{
		Description: "SSH Key resouce. Matches https://developers.bitlaunch.io/reference/ssh-key-object-1",

		CreateContext: resourceSSHKeyCreate,
		ReadContext:   resourceSSHKeyRead,
		DeleteContext: resourceSSHKeyDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the key.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"content": {
				Description: "The public portion of the SSH key.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"fingerprint": {
				Description: "The name of the key.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created": {
				Description: "The creation date of the key.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func setDataSSHKey(data *schema.ResourceData, key *gobitlaunch.SSHKey) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := data.Set("name", key.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("content", key.Content); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("fingerprint", key.Fingerprint); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("created", key.Created.Format(time.RFC3339)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceSSHKeyCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*apiClient).client
	tflog.Trace(ctx, "Creating an sshKey")

	key := gobitlaunch.SSHKey{
		Name:    data.Get("name").(string),
		Content: data.Get("content").(string),
	}

	newKey, err := client.SSHKey.Create(&key)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(newKey.ID)
	setDataSSHKey(data, newKey)
	tflog.Trace(ctx, fmt.Sprintf("created SSH Key with fingerprint %s", newKey.Fingerprint))

	return diags
}

func resourceSSHKeyRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*apiClient).client
	tflog.Trace(ctx, "Reading an sshKey")

	keys, err := client.SSHKey.List()
	if err != nil {
		return diag.FromErr(err)
	}
	for _, key := range keys {
		if key.ID == data.Id() {
			setDataSSHKey(data, &key)
			return diags
		}
	}

	tflog.Trace(ctx, "Server not found")
	data.SetId("")
	return diags
}

func resourceSSHKeyDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*apiClient).client
	tflog.Trace(ctx, "Deleting an sshKey")

	err := client.SSHKey.Delete(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
