package onepassword

import (
	"errors"
	"os/exec"
	"sync"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account email address",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account password",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account secret key",
			},
			"subdomain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account subdomain: <subdomain>.1password.com",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"item": dataSourceItem(),
		},
		ResourcesMap:  map[string]*schema.Resource{},
		ConfigureFunc: createOnePassClient,
	}
}

// The interface{} return value of this function is the meta parameter
// that will be passed into all resource CRUD functions.
// ref: https://www.terraform.io/docs/plugins/provider.html#configurefunc
func createOnePassClient(d *schema.ResourceData) (*OnePassClient, error) {
	opPath, err := exec.LookPath("op")
	if err != nil {
		msg = `
		Could not find 1Password CLI. Please install.
		See https://support.1password.com/command-line/ for instructions.
		`
		return nil, errors.New(msg)
	}
	op := &OnePassClient{
		Email:     d.Get("email").(string),
		Password:  d.Get("password").(string),
		SecretKey: d.Get("secret_key").(string),
		Subdomain: d.Get("subdomain").(string),
		OpPath:    opPath,
		mutex:     &sync.Mutex,
	}
}
