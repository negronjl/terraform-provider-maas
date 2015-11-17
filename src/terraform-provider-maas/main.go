package main

import (
	"github.com/hashicorp/terraform/plugin"
	"maas"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: maas.Provider,
	})
}
