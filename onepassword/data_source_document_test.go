package onepassword

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testDataSourceDocumentConfig = `
provider "onepassword" {
	email = "test@testing.com"
	password = "test-password"
	secret_key = "test-secret-key"
	subdomain = "test"
}

data "onepassword_document" "test" {
	vault = "test-vault"
	document = "test-doc"
}

output "test_doc" {
	value = "${file("${data.onepassword_document.test.path}")}"
} 
`

func TestDataSourceDocument(t *testing.T) {
	progPath, err := buildMockOnePassword()
	fmt.Println("progPath: ", progPath)
	if err != nil {
		t.Errorf("failed to build mock 1Password cli: %s", err)
	}

	os.Setenv("OP_PATH", progPath)

	docDir, err := ioutil.TempDir("", "documents")

	if err != nil {
		t.Errorf("error creating documents dir: %s", err)
	}

	docPath := filepath.Join(docDir, "test-doc")
	data := []byte("hello world")

	if err := ioutil.WriteFile(docPath, data, 0644); err != nil {
		t.Errorf("error creating mock file: %s", err)
	}

	defer os.RemoveAll(docDir)

	resource.UnitTest(t, resource.TestCase{
		Providers: createTestProviders(docDir),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDocumentConfig,
				Check: func(s *terraform.State) error {
					_, ok := s.RootModule().Resources["data.onepassword_document.test"]
					if !ok {
						return fmt.Errorf("missing data.onepassword_document.test data source")
					}

					outputs := s.RootModule().Outputs

					if outputs["test_doc"] == nil {
						return fmt.Errorf("missing 'test_doc' output")
					}

					if outputs["test_doc"].Value != "hello world" {
						return fmt.Errorf("'test_doc' != 'hello world'")
					}

					return nil
				},
			},
		},
	})

}
