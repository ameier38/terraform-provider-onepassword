package onepassword

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"path/filepath"
)

type vaultName string
type itemName string
type documentDir string
type documentName string
type documentPath string
type sectionName string
type fieldName string
type fieldValue string
type itemResponse []byte
type itemMap map[sectionName]map[fieldName]fieldValue

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

type documentGetter interface {
	getDocument(vaultName, itemName, documentDir)
}

type authenticator interface {
	authenticate() error
}

// Calls the `op signin` command and passes in the password via stdin.
// usage: op signin <signinaddress> <emailaddress> <secretkey> [--output=raw]
func (op *Client) authenticate() error {
	cmd := exec.Command(op.OpPath, "signin", op.Subdomain, op.Email, op.SecretKey, "--output=raw")
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
	op.Session = string(output)
	return nil
}

func getArg(key string, value string) string {
	return fmt.Sprintf("--%s=%s", key, value)
}

// Calls `op get item` command.
// usage: op get item <item> [--vault=<vault>] [--include-trash]
func (op Client) getItem(vault vaultName, item itemName) (itemResponse, error) {
	sessionArg := getArg("session", op.Session)
	vaultArg := getArg("vault", string(vault))
	cmd := exec.Command(string(op.OpPath), "get", "item", string(item), sessionArg, vaultArg)
	res, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("error calling 1Password: %s", err)
		return itemResponse(""), err
	}
	return itemResponse(res), nil
}

// Calls `op get document` command
// usage: op get document <document> [--vault=<vault>] > <filename>
func (op Client) getDocument(vault vaultName, docName documentName, docDir documentDir) (documentPath, error) {
	sessionArg := getArg("session", op.Session)
	vaultArg := getArg("vault", string(vault))
	docPath := filepath.Join(string(docDir), string(docName))
	cmd := exec.Command(string(op.OpPath), "get", "document", string(docName), docPath, sessionArg, vaultArg)
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("error calling 1Password: %s", err)
		return documentPath(""), err
	}
	return documentPath(docPath), nil
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
