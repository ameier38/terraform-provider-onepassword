package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/ameier38/terraform-provider-onepassword/onepassword"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	docDir, err := ioutil.TempDir("", "documents")

	if err != nil {
		log.Fatal("error creating documents directory")
		os.Exit(1)
	}

	defer os.RemoveAll(docDir)

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return onepassword.Provider(docDir)
		},
	})
}
