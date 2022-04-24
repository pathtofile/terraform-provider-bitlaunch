package tf_bitlaunch

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceBitlaunchServer(t *testing.T) {
	t.Skip("resource not yet implemented, remove this once you add your own code")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBitlaunchServer,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"bitlaunch_sshkey.sshkey", "name", regexp.MustCompile("^ba")),
				),
			},
		},
	})
}

const testAccResourceBitlaunchServer = `
provider "bitlaunch" {
	token = "aaa"
  }

resource "bitlaunch_sshkey" "sshkey" {
	name    = "tf_sshkeys"
	content = "aaaa"
  }  
`
