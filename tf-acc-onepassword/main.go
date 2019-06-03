package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const mockResponse = `
{
  "uuid": "test-item",
  "templateUuid": "102",
  "trashed": "N",
  "createdAt": "2019-05-18T14:58:54Z",
  "updatedAt": "2019-05-18T15:04:56Z",
  "itemVersion": 2,
  "vaultUuid": "test-vault",
  "details": {
    "fields": [],
    "notesPlain": "",
    "sections": [
      {
        "fields": [
          {
            "k": "menu",
            "n": "database_type",
            "t": "type",
            "v": "postgresql"
          },
          {
            "inputTraits": {
              "keyboard": "URL"
            },
            "k": "string",
            "n": "hostname",
            "t": "server",
            "v": "redshift.company.io"
          },
          {
            "inputTraits": {
              "keyboard": "NumberPad"
            },
            "k": "string",
            "n": "port",
            "t": "port",
            "v": "5439"
          },
          {
            "inputTraits": {
              "autocapitalization": "none",
              "autocorrection": "no"
            },
            "k": "string",
            "n": "database",
            "t": "database",
            "v": "test-db"
          },
          {
            "inputTraits": {
              "autocapitalization": "none",
              "autocorrection": "no"
            },
            "k": "string",
            "n": "username",
            "t": "username",
            "v": "test-user"
          },
          {
            "k": "concealed",
            "n": "password",
            "t": "password",
            "v": "test-password"
          },
          {
            "k": "string",
            "n": "sid",
            "t": "SID",
            "v": ""
          },
          {
            "k": "string",
            "n": "alias",
            "t": "alias",
            "v": ""
          },
          {
            "k": "string",
            "n": "options",
            "t": "connection options",
            "v": ""
          },
          {
            "k": "string",
            "n": "custom",
            "t": "schema",
            "v": "development"
          }
        ],
        "name": "",
        "title": ""
      }
    ]
  },
  "overview": {
    "URLs": [],
    "ainfo": "redshift.company.io",
    "pbe": 0,
    "pgrng": false,
    "ps": 0,
    "tags": [],
    "title": "Redshift",
    "url": ""
  }
}
`

func main() {
	input := bufio.NewReader(os.Stdin)
	// fake reading in password from stdin
	if _, err := input.ReadBytes('\n'); err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr, "error reading input:", err)
			os.Exit(1)
		}
	}
	fmt.Println(mockResponse)
	os.Exit(0)
}
