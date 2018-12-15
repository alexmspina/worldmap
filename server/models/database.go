package models

import (
	"fmt"

	"github.com/alexmspina/worldmap/server/helpers"
	"github.com/boltdb/bolt"
)

// DB main bolt database object
var DB, _ = SetupDB()

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
		_, err = root.CreateBucketIfNotExists([]byte("FLEET"))
		if err != nil {
			return fmt.Errorf("could not create fleet bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte("SATPOS"))
		if err != nil {
			return fmt.Errorf("could not create satellite positions bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte("CATSEYES"))
		if err != nil {
			return fmt.Errorf("could not create catseyes bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not setup buckets, %v", err)
	}
	fmt.Println("DB setup Done")
	return db, nil
}

// GetDbBucket pulls the desired bucket from the given database
func GetDbBucket(db *bolt.DB, mb string, b string, l *[][][]byte) {
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(mb)).Bucket([]byte(b))
		b.ForEach(func(k, v []byte) error {
			tmp := [][]byte{k, v}
			*l = append(*l, tmp)
			return nil
		})
		return nil
	})
	helpers.PanicErrors(err)
}

// GetDBObject retrieves the object from the desired bucket within a particular database and prints it to the screen
func GetDBObject(thing string, db *bolt.DB, dbBucket string, subBucket string) []byte {
	var b []byte
	db.View(func(tx *bolt.Tx) error {
		b = tx.Bucket([]byte(dbBucket)).Bucket([]byte(subBucket)).Get([]byte(thing))
		return nil
	})

	return b
}
