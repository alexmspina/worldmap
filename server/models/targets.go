package models

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

type featureCollection struct {
	Type     string    `json:"type"`
	Features []feature `json:"features"`
}

type feature struct {
	Type       string   `json:"type"`
	Geometry   geometry `json:"geometry"`
	Properties property `json:"properties"`
}

type geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type property struct {
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
func FillTargetsBucket(f string, db *bolt.DB, t time.Time) error {
	// Open file from filename
	r := OpenCSV(f)

	// Get column header values
	header := getHeader(r)

	// Create list of feature structs
	features := make([]feature, 0)
	createFeatureList(r, header, &features)

	// Create feature collection struct mirroring geojson collection
	fc := featureCollection{
		Type:     "featureCollection",
		Features: features,
	}

	featuresBytes, err := json.MarshalIndent(fc, "", "\t")
	if err != nil {
		panic(err)
	}
	featuresBytes = bytes.Replace(featuresBytes, []byte("\\u0026"), []byte("&"), -1)
	featuresBytes = bytes.Trim(featuresBytes, "\r")

	err = db.Update(func(tx *bolt.Tx) error {
		err = tx.Bucket([]byte("DB")).Bucket([]byte("TARGETS")).Put([]byte(t.Format(time.RFC3339)), featuresBytes)
		if err != nil {
			return fmt.Errorf("could not fill targets bucket: %v", err)
		}
		return nil
	})
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

func createFeatureList(r *csv.Reader, h []string, fl *[]feature) {
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if record[0] == h[0] {
			continue
		}

		f := buildFeature(record)
		*fl = append(*fl, f)
	}
}

func buildFeature(r []string) feature {
	lat := convertStringToFloat64(r[2])
	lon := convertStringToFloat64(r[3])
	if lon > 180.0 {
		lon = lon - 360.0
	}
	coordinates := []float64{lon, lat}
	geopoint := geometry{"Point", coordinates}
	prop := property{
		TargetID:    r[0],
		ShortName:   r[1],
		Altitude:    r[4],
		GatewayFlag: r[5],
		TTCFlag:     r[6],
		MinElTlmAOS: convertStringToFloat64(r[7]),
		MinElTlmLOS: convertStringToFloat64(r[8]),
		LongName:    r[9],
		FileCode:    r[10],
	}
	f := feature{
		"Feature",
		geopoint,
		prop,
	}
	return f
}

func convertStringToFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
