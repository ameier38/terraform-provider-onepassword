package onepassword

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDocument(docDir string) *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDocumentRead(docDir),

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
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDocumentRead(docDir string) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		op := meta.(*Client)
		vault := vaultName(d.Get("vault").(string))
		docName := documentName(d.Get("document").(string))
		docPath, err := op.getDocument(vault, docName, documentDir(docDir))
		if err != nil {
			return err
		}
		d.Set("path", string(docPath))
		d.SetId(time.Now().UTC().String())
		return nil
	}
}
