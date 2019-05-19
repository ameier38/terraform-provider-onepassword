package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

type VaultID string
type ItemID string
type Field string
type OPItemResponse string

type OPItemRequest struct {
	VaultID VaultID
	ItemID  ItemID
	Field   Field
}

type OPItem struct {
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

type IOPClient interface {
	GetItem(OPItemRequest) (OPItemResponse, error)
}

type OPClient struct{}

func (op OPClient) GetItem(vaultID VaultID, itemID ItemID, field Field) {
	vaultArg := fmt.Sprintf("--vault=%s", vaultID)
	opCmd := exec.Command("op", "get", "item", vaultArg, itemID)
	opOut, err := opCmd.Output()
	if err != nil {
		fmt.Println("error executing op command:", err)
		os.Exit(1)
	}
	parsedItem := opItem{}
	if err := json.Unmarshal(opOut, &parsedItem); err != nil {
		fmt.Println("could not unmarshal item:", string(opOut))
		os.Exit(1)
	}
	return parsedItem
}

func validateInput(vaultID string, itemID string, field string) (VaultID, ItemID, Field) {
	if vaultID == "" {
		fmt.Println("You must provide a vault UUID using -vault flag.")
	}
	if itemID == "" {
		fmt.Println("You must provide a item UUID using -item flag.")
	}
	if field == "" {
		fmt.Println("You must provide a field key using -field flag.")
	}
	return vaultID, itemID, field
}

func parseItem(response OPItemResponse) OPItem {
}

func getValue(item opItem, field string) string {
	for _, field := range item.Details.Sections[0].Fields {
		if field.Key == field {
			return field.Value
		}
	}
	return ""
}

func main() {
	vaultIDPtr := flag.String("vault", "", "The vault UUID")
	itemIDPtr := flag.String("item", "", "The item UUID")
	fieldPtr := flag.String("field", "", "The key from which to get the value")
	flag.Parse()
	vaultID, itemID, field := validateInput(*vaultIDPtr, *itemIDPtr, *fieldPtr)
	item := getItem(vaultID, itemID)
	value := getValue(item, key)
	if value == "" {
		fmt.Println("Could not get value.", "Have you logged in with 'iex $(op signin moneylionfinance)'?")
		os.Exit(1)
	}
	fmt.Printf(`{"value":"%s"}`, value)
}
