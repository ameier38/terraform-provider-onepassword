package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type vaultName string
type itemName string
type fieldName string
type itemResponse []byte
type fieldMap map[string]string
type response string

type itemRequest struct {
	VaultName vaultName `json:"vaultName"`
	ItemName  itemName  `json:"itemName"`
}

type parsedItem struct {
	UUID    string `json:"uuid"`
	Details struct {
		Sections []struct {
			Fields []struct {
				Key   string `json:"t"`
				Value string `json:"v"`
			} `json:"fields"`
		} `json:"sections"`
	} `json:"details"`
}

type itemGetter interface {
	getItem(itemRequest) (itemResponse, error)
}

type onePasswordItemGetter struct{}

func (onePasswordItemGetter) getItem(itemReq itemRequest) (itemResponse, error) {
	vaultArg := fmt.Sprintf("-vault=%s", itemReq.VaultName)
	itemArg := string(itemReq.ItemName)
	opCmd := exec.Command("op", "get", "item", vaultArg, itemArg)
	itemRes, err := opCmd.Output()
	if err != nil {
		err = fmt.Errorf("error calling 1Password; make sure you log in with 'iex $(op signin)': %s", err)
		return itemResponse(""), err
	}
	return itemRes, nil
}

func (res itemResponse) parse() (parsedItem, error) {
	pItem := parsedItem{}
	if err := json.Unmarshal(res, &pItem); err != nil {
		return pItem, err
	}
	return pItem, nil
}

func (pItem parsedItem) createFieldMap() fieldMap {
	m := make(fieldMap)
	for _, field := range pItem.Details.Sections[0].Fields {
		m[field.Key] = field.Value
	}
	return m
}

func getRequest(input []byte) (itemRequest, error) {
	itemReq := itemRequest{}
	if err := json.Unmarshal(input, &itemReq); err != nil {
		err = fmt.Errorf("error unmarshaling request: %s\n%s", string(input), err)
		return itemReq, err
	}
	return itemReq, nil
}

func getResponse(getter itemGetter, itemReq itemRequest) (response, error) {
	emtpyRes := response("")
	itemRes, err := getter.getItem(itemReq)
	if err != nil {
		return emtpyRes, err
	}
	pItem, err := itemRes.parse()
	if err != nil {
		return emtpyRes, err
	}
	m := pItem.createFieldMap()
	if err != nil {
		return emtpyRes, err
	}
	resData, err := json.Marshal(m)
	if err != nil {
		return emtpyRes, err
	}
	res := response(string(resData))
	return res, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadBytes('\n')
	if err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr, "error reading input:", err)
			os.Exit(1)
		}
	}
	itemReq, err := getRequest(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid request:", err)
		os.Exit(1)
	}
	getter := onePasswordItemGetter{}
	res, err := getResponse(getter, itemReq)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error!:", err)
		os.Exit(1)
	} else {
		fmt.Print(res)
	}
}
