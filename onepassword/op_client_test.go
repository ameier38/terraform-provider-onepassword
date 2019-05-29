package onepassword

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockOnePassClient struct {
	Session string
}

func (op *MockOnePassClient) authenticate() error {
	op.Session = "test-session"
}

func (MockOnePassClient) getItem(vault vaultName, item itemName) (itemResponse, error) {
	dat, err := ioutil.ReadFile("./testItemResponse.json")
	if err != nil {
		return itemResponse(""), err
	}
	return itemResponse(dat), nil
}

func TestGetItem(t *testing.T) {
	opClient := MockOnePassClient{}
	expectedItemMap := itemMap{
		sectionName(""): map[fieldName]fieldValue{
			fieldName("SID"):                fieldValue(""),
			fieldName("alias"):              fieldValue(""),
			fieldName("connection options"): fieldValue(""),
			fieldName("database"):           fieldValue("test-db"),
			fieldName("password"):           fieldValue("test-password"),
			fieldName("port"):               fieldValue("5439"),
			fieldName("schema"):             fieldValue("development"),
			fieldName("server"):             fieldValue("redshift.company.io"),
			fieldName("type"):               fieldValue("postgresql"),
			fieldName("username"):           fieldValue("test-user"),
		},
	}
	mockResponse, err := opClient.getItem(vaultName("test-vault"), itemName("test-item"))
	assert.Nil(t, err)
	actualItemMap, err := mockResponse.parse()
	if assert.Nil(t, err) {
		assert.Equal(t, expectedItemMap, actualItemMap, "expectedItemMap should equal actualItemMap")
	}
}
