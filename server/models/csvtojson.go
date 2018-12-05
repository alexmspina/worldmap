package models

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// COMMAND LINE INPUTS USING FLAGS
	path := flag.String("path", "nothing", "path to the csv file")
	flag.Parse()

	// DETERMINE OUTPUT DIRECTORY AND FILE TYPE FROM *PATH
	outputDir := filepath.Dir(*path) + "/" + filepath.Base(*path)[0:len(filepath.Base(*path))-4] + ".json"
	x := filepath.Ext(*path)

	// OPEN FILE FROM PATH
	f, err := os.Open(*path)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// GENERATE JSON BASED ON FILE EXTENSION
	if x == ".csv" {
		r := csv.NewReader(f)
		a, err := r.ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		p := a[0]

		jsonMap := make([]map[string]string, 0)
		for i := 1; i < len(a); i++ {
			data := a[i]
			objectMap := make(map[string]string)
			for j, key := range p {
				objectMap[key] = data[j]
			}
			jsonMap = append(jsonMap, objectMap)
		}

		rawJSON := formatJSON(jsonMap)
		writeJSONFile(rawJSON, outputDir)

	} else if x == ".tle" {
		r := io.Reader(f)
		a, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}
		s := string(a)
		l := strings.Split(s, "\n")
		if len(l[len(l)-1]) <= 0 {
			l = l[:len(l)-1]
		}
		if len(l)%3 != 0 {
			l = l[1:]
		}

		jsonMap := make([]map[string]string, 0)
		for i := 0; i < len(l)/3; i++ {
			objectMap := make(map[string]string)
			objectMap["name"] = strings.Trim(l[i*3], "\r")
			objectMap["line1"] = strings.Trim(l[i*3+1], "\r")
			objectMap["line2"] = strings.Trim(l[i*3+2], "\r")
			jsonMap = append(jsonMap, objectMap)
		}

		rawJSON := formatJSON(jsonMap)
		writeJSONFile(rawJSON, outputDir)
	}
}

// JSON FORMATTING
func formatJSON(j []map[string]string) []byte {
	b, err := json.MarshalIndent(j, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	b = bytes.Trim(b, "\r")
	return b
}

func writeJSONFile(b []byte, f string) {
	err := ioutil.WriteFile(f, b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
