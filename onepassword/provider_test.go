package onepassword

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func createTestProviders() map[string]terraform.ResourceProvider {
	testProvider := Provider()
	testProviders := map[string]terraform.ResourceProvider{
		"onepassword": testProvider,
	}
	return testProviders
}
