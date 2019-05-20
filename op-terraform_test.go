package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockItemGetter struct{}

func (mockItemGetter) getItem(itemRequest) (itemResponse, error) {
	dat, err := ioutil.ReadFile("./testItemResponse.json")
	if err != nil {
		return itemResponse(""), err
	}
	return itemResponse(dat), nil
}

func TestGetRequest(t *testing.T) {
	input := []byte(`{"vaultName": "test-vault", "itemName": "test-item"}`)
	actualReq, err := getRequest(input)
	expectedReq := itemRequest{vaultName("test-vault"), itemName("test-item")}
	if assert.Nil(t, err) {
		assert.Equal(t, expectedReq, actualReq, "expectedReq should equal actualReq")
	}
}

func TestGetResponse(t *testing.T) {
	expectedResponse := `{
		"SID": "",
		"alias": "",
		"connection options": "",
		"database": "test-db",
		"password": "test-password",
		"port": "5439",
		"schema": "development",
		"server": "redshift.company.io",
		"type": "postgresql",
		"username": "test-user"
	}`
	getter := mockItemGetter{}
	itemReq := itemRequest{vaultName("test-vault"), itemName("test-item")}
	actualResponse, err := getResponse(getter, itemReq)
	if assert.Nil(t, err) {
		assert.JSONEq(t, expectedResponse, string(actualResponse), "actualResponse should equal expectedResponse")
	}
}
