package tf_bitlaunch

import (
	"context"
	"fmt"

	"github.com/bitlaunchio/gobitlaunch"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// https://developers.bitlaunch.io/reference/view-host-create-options
var HostIDs = map[string]int{
	"DigitalOcean": 0,
	"Vultr":        1,
	"Linode":       2,
	"BitLaunch":    4,
}

func ValidateHostID(val interface{}, key string) (warns []string, errs []error) {
	hostName := val.(string)
	hostNames := maps.Keys(HostIDs)
	if !slices.Contains(hostNames, hostName) {
		errs = append(errs, fmt.Errorf("%q Must be one of %s", key, hostNames))
	}
	return
}

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"token": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("BITLAUNCH_API_TOKEN", nil),
					Description: "API Token",
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"bitlaunch_sshkey": resourceSSHKey(),
				"bitlaunch_server": resourceServer(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"bitlaunch_size":   dataSourceSize(),
				"bitlaunch_region": dataSourceRegion(),
				"bitlaunch_image":  dataSourceImage(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	client *gobitlaunch.Client
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
		token := data.Get("token").(string)

		client := gobitlaunch.NewClient(token)

		return &apiClient{client: client}, nil
	}
}
