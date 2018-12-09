package models

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

// GetHeader from csv files
func GetHeader(cr *csv.Reader) []string {
	record, err := cr.Read()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		log.Fatal(err)
	}

	return record
}

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
