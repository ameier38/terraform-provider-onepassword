package onepassword

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testDataSourceItemConfig = `
provider "onepassword" {
	email = "test@testing.com"
	password = "test-password"
	secret_key = "test-secret-key"
	subdomain = "test"
}

data "onepassword_item" "test" {
	vault = "test-vault"
	item = "test-item"
	section = "Terraform"
}

output "test_user" {
	value = "${data.onepassword_item.test.result["username"]}"
} 

output "test_password" {
	value = "${data.onepassword_item.test.result["password"]}"
}
`

func TestDataSourceItem(t *testing.T) {
	progPath, err := buildMockOnePassword()
	if err != nil {
		t.Errorf("failed to build mock 1Password cli: %s", err)
	}

	os.Setenv("OP_PATH", progPath)

	resource.UnitTest(t, resource.TestCase{
		Providers: createTestProviders(),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceItemConfig,
				Check: func(s *terraform.State) error {
					_, ok := s.RootModule().Resources["data.onepassword_item.test"]
					if !ok {
						return fmt.Errorf("missing data.onepassword_item.test data source")
					}

					outputs := s.RootModule().Outputs

					if outputs["test_user"] == nil {
						return fmt.Errorf("missing 'test_user' output")
					}

					if outputs["test_password"] == nil {
						return fmt.Errorf("missing 'test_password' output")
					}

					if outputs["test_user"].Value != "test-user" {
						return fmt.Errorf("'%s' != 'test-user'", outputs["test_user"].Value)
					}

					if outputs["test_password"].Value != "test-password" {
						return fmt.Errorf("'%s' != 'test-password'", outputs["test_password"].Value)
					}

					return nil
				},
			},
		},
	})

}
