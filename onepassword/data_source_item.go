package onepassword

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceItem() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRead(),

		Schema: map[string]*schema.Schema{
			"vault": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "1Password Vault in which item resides",
			},
			"section": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Section in which field resides",
			},
			"field": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Field for which to get the value",
			},
			"result": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceItemRead(d *schema.ResourceData, op *OnePassClient) error {
	vault := d.Get("vault").(string)
	section := d.Get("section").(string)
	field := d.Get("field").(string)
}
