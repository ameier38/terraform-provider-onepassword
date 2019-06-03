package onepassword

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
}

output "test_user" {
	value = "${data.onepassword_item.test.result["username"]}"
} 

output "test_password" {
	value = "${data.onepassword_item.test.result["password"]}"
}
`

func getExtension() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

func buildMockOnePassword() (string, error) {
	cmd := exec.Command(
		"go",
		"install",
		"github.com/ameier38/terraform-provider-onepassword/tf-acc-onepassword")

	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build mock op program: %s\n%s", err, output)
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(os.Getenv("HOME"), "go")
	}

	programPath := filepath.Join(
		filepath.SplitList(gopath)[0],
		"bin",
		"tf-acc-onepassword"+getExtension())

	return programPath, nil
}

func TestDataSourceItem(t *testing.T) {
	progPath, err := buildMockOnePassword()
	if err != nil {
		t.Errorf("failed to build mock 1Password cli: %s", err)
	}

	os.Setenv("OP_PATH", progPath)

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
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
						return fmt.Errorf("'test_user' != 'test-user'")
					}

					if outputs["test_password"].Value != "test-password" {
						return fmt.Errorf("'test_password' != 'test-password'")
					}

					return nil
				},
			},
		},
	})

}
