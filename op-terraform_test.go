package main

import (
	"io/ioutil"
	"testing"
)

type mockItemGetter struct{}

func (mockItemGetter) getItem(itemRequest) (itemResponse, error) {
	dat, err := ioutil.ReadFile("./testItemResponse.json")
	if err != nil {
		return itemResponse(""), err
	}
	return itemResponse(dat), nil
}

func TestGetFieldValueResponse(t *testing.T) {
	tables := []struct {
		fieldName          string
		fieldValueResponse string
	}{
		{"server", `{"value":"redshift.company.io"}`},
		{"username", `{"value":"test-user"}`},
		{"password", `{"value":"test-password"}`},
		{"schema", `{"value":"development"}`},
	}
	for _, table := range tables {
		itemReq := itemRequest{vaultID("test-vault"), itemID("test-item"), fieldName(table.fieldName)}
		getter := mockItemGetter{}
		fieldValueRes, err := getFieldValueResponse(getter, itemReq)
		if err != nil {
			t.Fatalf("Expected err to be nil but got %s", err)
		}
		if string(fieldValueRes) != table.fieldValueResponse {
			t.Errorf(
				"expectedFieldValueResponse: %s did not match actualFieldValueResponse: %s",
				string(table.fieldValueResponse),
				string(fieldValueRes))
		}
	}
}
