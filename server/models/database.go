package models

import (
	"fmt"

	"github.com/boltdb/bolt"
)

// SetupDB initializes a bolt database
func SetupDB() (*bolt.DB, error) {
	db, err := bolt.Open("system.db", 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("DB"))
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte("TARGETS"))
		if err != nil {
			return fmt.Errorf("could not create targets bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte("BEAMPLANS"))
		if err != nil {
			return fmt.Errorf("could not create beamplans bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte("ZONES"))
		if err != nil {
			return fmt.Errorf("could not create zones bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte("TLES"))
		if err != nil {
			return fmt.Errorf("could not create tle bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not setup buckets, %v", err)
	}
	fmt.Println("DB setup Done")
	return db, nil
}
