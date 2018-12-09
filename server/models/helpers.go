package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
)

// FormatJSON FORMATTING
func FormatJSON(j []map[string]string) []byte {
	b, err := json.MarshalIndent(j, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	b = bytes.Trim(b, "\r")
	return b
}

// ConvertStringToFloat64 converts with error check
func ConvertStringToFloat64(s string) float64 {
	switch "" {
	case s:
		return 0
	default:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
		return f
	}
}

// ConvertFloat64ToString takes a float and returns a string
func ConvertFloat64ToString(f float64) string {
	return fmt.Sprintf("%f", f)
}

// Degs2Rads converts float64 values from degrees to radians
func Degs2Rads(d float64) float64 {
	r := d * math.Pi / 180
	return r
}

// Rads2Degs converts float64 values from radians to degrees
func Rads2Degs(r float64) float64 {
	d := r * 180 / math.Pi
	return d
}

// PanicErrors process errors when they occur
func PanicErrors(err error) {
	if err != nil {
		panic(err)
	}
}

// StringInSlice determines if a string is in a given string slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
