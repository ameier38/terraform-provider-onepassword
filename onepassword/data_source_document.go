package onepassword

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDocument() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDocumentRead,

		Schema: map[string]*schema.Schema{
			"vault": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "1Password Vault in which document resides",
			},
			"document": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "1Password document to retrieve",
			},
			"result": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDocumentRead(d *schema.ResourceData, meta interface{}) error {
	op := meta.(*Client)
	vault := vaultName(d.Get("vault").(string))
	docName := documentName(d.Get("document").(string))
	docValue, err := op.getDocument(vault, docName)
	if err != nil {
		return err
	}
	d.Set("result", string(docValue))
	d.SetId(time.Now().UTC().String())
	return nil
}
