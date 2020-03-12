/* provider is a Terraform provider for MAAS.
 */
package provider

import (
	"github.com/hashicorp/terraform/terraform"
	"github.com/roblox/terraform-provider-maas/pkg/gmaw"

	"github.com/hashicorp/terraform/helper/schema"
)

// Provider exports a Terraform resource provider and is meant to be called from main().
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The MAAS API key",
			},
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The MAAS server URL (eg http://1.2.3.4:5240/MAAS)",
			},
			"api_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "2.0",
				Description: "The MAAS API version (default 2.0)",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"maas_instance":           resourceInstance(),
			"maas_interface_physical": ResourceNetworkInterfacePhysical(),
			"maas_interface_link":     ResourceNetworkInterfaceLink(),
			"maas_server":             ResourceServer(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"maas_subnet":          DataSubnet(),
			"maas_rack_controller": DataRackController(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return gmaw.GetClient(d.Get("api_url").(string), d.Get("api_version").(string), d.Get("api_key").(string))
}
