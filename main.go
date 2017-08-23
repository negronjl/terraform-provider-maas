package main

import (
	"github.com/hashicorp/terraform/plugin"
)

// Terraform plugin load point
func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: Provider,
	})
}
