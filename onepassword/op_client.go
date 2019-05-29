package onepassword

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

type subdomain string
type email string
type password string
type secretKey string
type session string
type opPath string
type vaultName string
type itemName string
type sectionName string
type fieldName string
type fieldValue string
type itemResponse []byte
type itemMap map[sectionName]map[fieldName]fieldValue
type response string

type OnePassClient struct {
	Subdomain subdomain
	Email     email
	Password  password
	SecretKey secretKey
	OpPath    opPath
	Session   session
}

type parsedItem struct {
	UUID    string `json:"uuid"`
	Details struct {
		Sections []struct {
			Name   string `json:"name"`
			Fields []struct {
				Key   string `json:"t"`
				Value string `json:"v"`
			} `json:"fields"`
		} `json:"sections"`
	} `json:"details"`
}

type itemGetter interface {
	getItem(vaultName, itemName) (itemResponse, error)
}

type authenticator interface {
	authenticate() error
}

// Calls the `op signin` command and passes in the password via stdin.
// usage: op signin <signinaddress> <emailaddress> <secretkey> [--output=raw]
func (op *OnePassClient) authenticate() error {
	cmd := exec.Command(
		string(op.OpPath),
		"signin",
		string(op.Subdomain),
		string(op.Email),
		string(op.SecretKey),
		"--output=raw")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("Cannot attach to stdin: %s", err)
	}
	go func() {
		defer stdin.Close()
		if _, err := io.WriteString(stdin, fmt.Sprintf("%s\n", op.Password)); err != nil {
			log.Println("[Error]", err)
		}
	}()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Cannot signin: %s", err)
	}
	op.Session = session(output)
	return nil
}

// Calls `op get item` command.
// usage: op get item <item> [--vault=<vault>] [--include-trash]
func (op OnePassClient) getItem(vault vaultName, item itemName) (itemResponse, error) {
	sessionArg := fmt.Sprintf("--session=%s", strings.Trim(string(op.Session), "\n"))
	vaultArg := fmt.Sprintf("--vault=%s", strings.Trim(string(vault), "\n"))
	cmd := exec.Command(string(op.OpPath), "get", "item", sessionArg, vaultArg, string(item))
	res, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("error calling 1Password: %s", err)
		return itemResponse(""), err
	}
	return itemResponse(res), nil
}

func (res itemResponse) parse() (itemMap, error) {
	im := make(itemMap)
	pItem := parsedItem{}
	if err := json.Unmarshal(res, &pItem); err != nil {
		return im, err
	}
	for _, section := range pItem.Details.Sections {
		fm := make(map[fieldName]fieldValue)
		for _, field := range section.Fields {
			fm[fieldName(field.Key)] = fieldValue(field.Value)
		}
		im[sectionName(section.Name)] = fm
	}
	return im, nil
}
