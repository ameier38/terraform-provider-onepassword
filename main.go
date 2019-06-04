package main

import (
	"github.com/ameier38/terraform-provider-onepassword/onepassword"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return onepassword.Provider()
		},
	})
}
