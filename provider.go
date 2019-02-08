package main

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider creates the schema for the provider config
func Provider() terraform.ResourceProvider {
	log.Println("[DEBUG] Initializing the MAAS provider")
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The api key for API operations",
			},
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The MAAS server URL. ie: http://1.2.3.4:80/MAAS",
			},
			"api_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "2.0",
				Description: "The MAAS API version. Currently: 1.0",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"maas_instance": resourceMAASMachine(), //deprecated, old name
			"maas_machine":  resourceMAASMachine(),
		},

		ConfigureFunc: providerConfigure,
	}
}

// providerConfigure loads in the provider configuration
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Println("[DEBUG] Configuring the MAAS provider")
	config := Config{
		APIKey: d.Get("api_key").(string),
		APIURL: d.Get("api_url").(string),
		APIver: d.Get("api_version").(string),
	}
	return config.Client()
}
