package onepassword

import (
	"fmt"
	"sync"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider : 1Password Provider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"op": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OP_PATH", "op"),
				Description: "Path to 1Password CLI (i.e, 'op'). Defaults to value of OP_PATH env variable.",
			},
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
			"onepassword_item":     dataSourceItem(),
			"onepassword_document": dataSourceDocument(),
		},
		ResourcesMap:  map[string]*schema.Resource{},
		ConfigureFunc: createClient,
	}
}

// The interface{} return value of this function is the meta parameter
// that will be passed into all resource CRUD functions.
// ref: https://www.terraform.io/docs/plugins/provider.html#configurefunc
func createClient(d *schema.ResourceData) (interface{}, error) {
	op := &Client{
		OpPath:    d.Get("op").(string),
		Email:     d.Get("email").(string),
		Password:  d.Get("password").(string),
		SecretKey: d.Get("secret_key").(string),
		Subdomain: d.Get("subdomain").(string),
		mutex:     &sync.Mutex{},
	}
	if err := op.authenticate(); err != nil {
		return nil, fmt.Errorf("could not authenticate 1Password: %s", err)
	}
	return op, nil
}
