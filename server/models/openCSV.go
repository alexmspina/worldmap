package models

import (
	"encoding/csv"
	"fmt"
	"os"
)

// OpenCSV opens a csv file
func OpenCSV(f string) *csv.Reader {
	file, err := os.Open(f)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	r := csv.NewReader(file)

	return r
}
