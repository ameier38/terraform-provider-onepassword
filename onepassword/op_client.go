package onepassword

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

type vaultName string
type itemName string
type documentName string
type documentValue string
type sectionName string
type fieldName string
type fieldValue string
type itemResponse []byte
type fieldMap map[fieldName]fieldValue
type sectionMap map[sectionName]fieldMap

// Client : 1Password client
type Client struct {
	OpPath    string
	Subdomain string
	Email     string
	Password  string
	SecretKey string
	Session   string
}

type parsedItem struct {
	UUID    string `json:"uuid"`
	Details struct {
		Sections []struct {
			Title  string `json:"title"`
			Fields []struct {
				Key   string `json:"t"`
				Value string `json:"v"`
			} `json:"fields"`
		} `json:"sections"`
	} `json:"details"`
}

// Calls the `op signin` command and passes in the password via stdin.
// usage: op signin <signinaddress> <emailaddress> <secretkey> [--output=raw]
func (op *Client) authenticate() error {
	signinAddress := fmt.Sprintf("%s.1password.com", op.Subdomain)
	cmd := exec.Command(op.OpPath, "signin", signinAddress, op.Email, op.SecretKey, "--output=raw")
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
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Cannot signin: %s\n%s", err, output)
	}
	op.Session = strings.Trim(string(output), "\n")
	return nil
}

func getArg(key string, value string) string {
	return fmt.Sprintf("--%s=%s", key, value)
}

func (res itemResponse) parseResponse() (sectionMap, error) {
	sm := make(sectionMap)
	pItem := parsedItem{}
	if err := json.Unmarshal(res, &pItem); err != nil {
		return sm, err
	}
	for _, section := range pItem.Details.Sections {
		fm := make(fieldMap)
		for _, field := range section.Fields {
			fm[fieldName(field.Key)] = fieldValue(field.Value)
		}
		sm[sectionName(section.Title)] = fm
	}
	return sm, nil
}

// Calls `op get item` command.
// usage: op get item <item> [--vault=<vault>] [--session=<session>]
func (op Client) getItem(vault vaultName, item itemName) (itemResponse, error) {
	sessionArg := getArg("session", op.Session)
	vaultArg := getArg("vault", string(vault))
	debugCmd := fmt.Sprintf("op get item %s %s %s", string(item), vaultArg, sessionArg)
	cmd := exec.Command(string(op.OpPath), "get", "item", string(item), vaultArg, sessionArg)
	res, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("error calling 1Password: %s\n%s\n'%s'", err, res, debugCmd)
		return itemResponse(""), err
	}
	return itemResponse(res), nil
}

// Calls `op get document` command
// usage: op get document <document> [--vault=<vault>] [--session=<session>]
func (op Client) getDocument(vault vaultName, docName documentName) (documentValue, error) {
	sessionArg := getArg("session", op.Session)
	vaultArg := getArg("vault", string(vault))
	debugCmd := fmt.Sprintf("op get document %s %s %s", string(docName), vaultArg, sessionArg)
	cmd := exec.Command(string(op.OpPath), "get", "document", string(docName), vaultArg, sessionArg)
	res, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("error calling 1Password: %s\n%s\n'%s'", err, res, debugCmd)
		return documentValue(""), err
	}
	return documentValue(res), nil
}
