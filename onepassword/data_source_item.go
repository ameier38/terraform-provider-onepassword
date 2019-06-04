package onepassword

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceItem() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceItemRead,

		Schema: map[string]*schema.Schema{
			"vault": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "1Password Vault in which item resides",
			},
			"item": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "1Password item to retrieve",
			},
			"section": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Section in which fields reside",
				Default:     "Terraform",
			},
			"result": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     schema.TypeString,
			},
		},
	}
}

func dataSourceItemRead(d *schema.ResourceData, meta interface{}) error {
	op := meta.(*Client)
	vault := vaultName(d.Get("vault").(string))
	item := itemName(d.Get("item").(string))
	section := sectionName(d.Get("section").(string))
	itemRes, err := op.getItem(vault, item)
	if err != nil {
		return err
	}
	sectionMap, err := itemRes.parseResponse()
	if err != nil {
		return err
	}
	result := sectionMap[section]
	d.Set("result", result)
	d.SetId(time.Now().UTC().String())
	return nil
}
