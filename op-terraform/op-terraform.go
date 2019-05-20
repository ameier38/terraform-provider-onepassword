package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type vaultID string
type itemID string
type fieldName string
type fieldValue string
type fieldValueResponse string
type itemResponse []byte

type itemRequest struct {
	vaultID   vaultID
	itemID    itemID
	fieldName fieldName
}

func (req itemRequest) validate() error {
	var err error
	var errStrings []string
	if req.vaultID == "" {
		errStrings = append(errStrings, "You must provide a vault UUID using -vault flag.")
	}
	if req.itemID == "" {
		errStrings = append(errStrings, "You must provide a item UUID using -item flag.")
	}
	if req.fieldName == "" {
		errStrings = append(errStrings, "You must provide a field key using -field flag.")
	}
	if len(errStrings) > 0 {
		err = fmt.Errorf(strings.Join(errStrings, "\n"))
	}
	return err
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
	vaultArg := fmt.Sprintf("--vault=%s", itemReq.vaultID)
	opCmd := exec.Command("op", "get", "item", vaultArg, string(itemReq.itemID))
	itemRes, err := opCmd.Output()
	if err != nil {
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

func (pItem parsedItem) fieldValue(fName fieldName) (fieldValue, error) {
	for _, field := range pItem.Details.Sections[0].Fields {
		if field.Key == string(fName) {
			return fieldValue(field.Value), nil
		}
	}
	return "", fmt.Errorf("could not find field %s", string(fName))
}

func getFieldValueResponse(getter itemGetter, itemReq itemRequest) (fieldValueResponse, error) {
	emtpyValueRes := fieldValueResponse("")
	itemRes, err := getter.getItem(itemReq)
	if err != nil {
		return emtpyValueRes, err
	}
	pItem, err := itemRes.parse()
	if err != nil {
		return emtpyValueRes, err
	}
	fValue, err := pItem.fieldValue(itemReq.fieldName)
	if err != nil {
		return emtpyValueRes, err
	}
	valueRes := fieldValueResponse(fmt.Sprintf(`{"value":"%s"}`, string(fValue)))
	return valueRes, nil
}

func main() {
	vaultIDPtr := flag.String("vault", "", "The vault UUID")
	itemIDPtr := flag.String("item", "", "The item UUID")
	fieldNamePtr := flag.String("field", "", "The key from which to get the value")
	flag.Parse()
	itemReq := itemRequest{vaultID(*vaultIDPtr), itemID(*itemIDPtr), fieldName(*fieldNamePtr)}
	if err := itemReq.validate(); err != nil {
		fmt.Fprintln(os.Stderr, "invalid request:", err)
	}
	getter := onePasswordItemGetter{}
	fieldValueRes, err := getFieldValueResponse(getter, itemReq)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error!:", err)
	} else {
		fmt.Print(fieldValueRes)
	}
}
