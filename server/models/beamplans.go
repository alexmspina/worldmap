package models

import (
	"encoding/csv"
	"fmt"
	"time"

	"github.com/alexmspina/worldmap/server/helpers"
	"github.com/boltdb/bolt"
)

// FillBeamplanBucket fills bolt db bucket with beamplan from initial files
func FillBeamplanBucket(f string, db *bolt.DB, t time.Time) error {
	r := OpenCSV(f)
	jsonMap := make([]map[string]string, 0)
	createBeamplanCollection(r, &jsonMap)
	rawJSON := helpers.FormatJSON(jsonMap)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("DB")).Bucket([]byte("BEAMPLANS")).Put([]byte(t.Format(time.RFC3339)), rawJSON)
		if err != nil {
			return fmt.Errorf("could not fill beamplans bucket: %v", err)
		}
		return nil
	})
	fmt.Println("Beamplans bucket filled")

	return err
}

func createBeamplanCollection(r *csv.Reader, j *[]map[string]string) {
	// write reader into array by line
	readerList, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	// get csv header values from first line
	h := readerList[0]

	// create a list of maps from data
	for i := 1; i < len(readerList); i++ {
		data := readerList[i]
		objectMap := make(map[string]string)
		for j, key := range h {
			objectMap[key] = data[j]
		}
		*j = append(*j, objectMap)
	}
}
