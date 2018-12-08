package models

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
)

type targetFeature struct {
	Type       string           `json:"type"`
	Geometry   PointGeometry    `json:"geometry"`
	Properties targetProperties `json:"properties"`
}

type targetProperties struct {
	TargetID    string  `json:"targetID"`
	ShortName   string  `json:"shortName"`
	Altitude    string  `json:"altitude"`
	GatewayFlag string  `json:"gatewayFlag"`
	TTCFlag     string  `json:"ttcFlag"`
	MinElTlmAOS float64 `json:"minElTlmAOS"`
	MinElTlmLOS float64 `json:"minElTlmLOS"`
	LongName    string  `json:"longName"`
	FileCode    string  `json:"fileCode"`
}

// FillTargetsBucket initializes a BoltDB bucket from TARGETS file
func FillTargetsBucket(f string, db *bolt.DB) error {
	// Open file from filename
	r := OpenCSV(f)

	// Get column header values
	header := getHeader(r)
	var err error

	// Create list of feature structs
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if record[0] == header[0] {
			continue
		}

		f := buildTargetFeature(record)
		id := f.Properties.TargetID

		featureBytes, err := json.MarshalIndent(f, "", "\t")
		if err != nil {
			panic(err)
		}
		featureBytes = bytes.Replace(featureBytes, []byte("\\u0026"), []byte("&"), -1)
		featureBytes = bytes.Trim(featureBytes, "\r")

		err = db.Update(func(tx *bolt.Tx) error {
			err = tx.Bucket([]byte("DB")).Bucket([]byte("TARGETS")).Put([]byte(id), featureBytes)
			if err != nil {
				return fmt.Errorf("could not fill targets bucket: %v", err)
			}
			return nil
		})
	}
	fmt.Println("Targets bucket filled")

	return err
}

func getHeader(cr *csv.Reader) []string {
	record, err := cr.Read()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		log.Fatal(err)
	}

	return record
}

func buildTargetFeature(r []string) targetFeature {
	lat := ConvertStringToFloat64(r[2])
	lon := ConvertStringToFloat64(r[3])
	if lon > 180.0 {
		lon = lon - 360.0
	}
	coordinates := []float64{lon, lat}
	geopoint := PointGeometry{"Point", coordinates}
	props := targetProperties{
		TargetID:    r[0],
		ShortName:   r[1],
		Altitude:    r[4],
		GatewayFlag: r[5],
		TTCFlag:     r[6],
		MinElTlmAOS: ConvertStringToFloat64(r[7]),
		MinElTlmLOS: ConvertStringToFloat64(r[8]),
		LongName:    r[9],
		FileCode:    r[10],
	}
	f := targetFeature{
		"Feature",
		geopoint,
		props,
	}
	return f
}

func convertStringToFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
