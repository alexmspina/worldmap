package models

import (
	"bytes"
	"encoding/json"
	"log"
)

// FormatJSON generic json formatter using marshal indent
func FormatJSON(j []map[string]string) []byte {
	b, err := json.MarshalIndent(j, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	b = bytes.Trim(b, "\r")
	return b
}
