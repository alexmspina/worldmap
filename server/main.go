package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/alexmspina/worldmap/server/models"
	"github.com/boltdb/bolt"
)

func main() {
	// parse command-line flag to determine root directory location of necessary files
	dir := flag.String("dir", "No directory provided", "input the directory where the initial files are located")
	flag.Parse()
	if *dir == "No directory provided" {
		fmt.Println("No directory provided for initial setup.")
		os.Exit(1)
	}

	// use Walk function to traverse root directory provided and create a list of files
	files := make([]string, 0)
	GetFilesFromDirectory(&files, *dir)

	// Create map of regular expressions
	regexmap := make(map[string]*regexp.Regexp, 0)
	preregexlist := []string{"TARGETS", "BEAMPLAN_LONGFORMAT", "ZONES", "ephemeris"}
	CreateRegexp(regexmap, preregexlist)

	// Setup bolt database
	db, _ := models.SetupDB()

	// Process the selected files depending on their type and fill bolt db buckets
	ProcessInitFiles(files, regexmap, db)

	// err := db.View(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte("DB")).Bucket([]byte("CATSEYES"))
	// 	// b.ForEach(func(k, v []byte) error {
	// 	// 	fmt.Println(string(k), string(v))
	// 	// 	return nil
	// 	// })
	// 	// return nil
	// 	v := b.Get([]byte("1"))
	// 	fmt.Println(string(v))
	// 	return nil
	// })
	// models.PanicErrors(err)
}

// ProcessInitFiles checks file names against list of regular expressions and calls handlers based on results
func ProcessInitFiles(files []string, regexmap map[string]*regexp.Regexp, db *bolt.DB) {
	for _, file := range files {
		switch true {
		case regexmap["TARGETS"].MatchString(filepath.Base(file)):
			models.FillTargetsBucket(file, db)
		case regexmap["BEAMPLAN_LONGFORMAT"].MatchString(filepath.Base(file)):
			// models.FillBeamplanBucket(file, db, time.Now())
			fmt.Printf("Beamplan file found: %s\n\n", file)
		case regexmap["ZONES"].MatchString(filepath.Base(file)):
			models.FillZonesBucket(file, db)
			models.FillCatseyesBucket(db)
		case regexmap["ephemeris"].MatchString(filepath.Base(file)):
			fmt.Printf("TLE file found: %s\n\n", file)
		default:
			fmt.Printf("This file was not processed: %s\n\n", file)
		}
	}
}

// CreateRegexp creates a map of string keys and their regular expression counterpart values
func CreateRegexp(r map[string]*regexp.Regexp, p []string) {
	for _, s := range p {
		regex, err := regexp.Compile(s)
		models.PanicErrors(err)
		r[s] = regex
	}
}

// GetFilesFromDirectory use the Walk function to create a list of files found in the given directory
func GetFilesFromDirectory(f *[]string, d string) {
	err := filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		*f = append(*f, path)
		return nil
	})
	models.PanicErrors(err)
}
