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

	bpregexmap := make(map[string]*regexp.Regexp, 0)
	bppreregex := []string{"mute", "B3", "M001", "M013"}
	CreateRegexp(bpregexmap, bppreregex)

	bpfilelist := make(map[string]string, 0)
	GetBeamplanFiles(files, regexmap["BEAMPLAN_LONGFORMAT"], bpregexmap, bpfilelist)

	// Setup bolt database
	db, _ := models.SetupDB()

	// Process the selected files depending on their type and fill bolt db buckets
	ProcessInitFiles(files, regexmap, db, bpfilelist)

	// err := db.View(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte("DB")).Bucket([]byte("FLEET"))
	// 	// b.ForEach(func(k, v []byte) error {
	// 	// 	fmt.Println(string(k), string(v))
	// 	// 	return nil
	// 	// })
	// 	// return nil
	// 	v := b.Get([]byte("M003"))
	// 	fmt.Println(string(v))
	// 	return nil
	// })
	// models.PanicErrors(err)
}

// GetBeamplanFiles organizes beamplan files into map sorted by active, spare, M001, M013
func GetBeamplanFiles(files []string, fileregex *regexp.Regexp, regexmap map[string]*regexp.Regexp, bpfiles map[string]string) {
	for _, file := range files {
		switch true {
		case fileregex.MatchString(filepath.Base(file)):
			switch true {
			case regexmap["mute"].MatchString(filepath.Base(file)):
				bpfiles["SPARE"] = file
			case regexmap["B3"].MatchString(filepath.Base(file)):
				bpfiles["B3"] = file
			case regexmap["M001"].MatchString(filepath.Base(file)):
				bpfiles["M001"] = file
			case regexmap["M013"].MatchString(filepath.Base(file)):
				bpfiles["M013"] = file
			default:
				bpfiles["ACTIVE"] = file
			}
		}
	}
}

// ProcessInitFiles checks file names against list of regular expressions and calls handlers based on results
func ProcessInitFiles(files []string, regexmap map[string]*regexp.Regexp, db *bolt.DB, bpfilelist map[string]string) {
	for _, file := range files {
		switch true {
		case regexmap["TARGETS"].MatchString(filepath.Base(file)):
			models.FillTargetsBucket(file, db)
		case regexmap["ZONES"].MatchString(filepath.Base(file)):
			models.FillZonesBucket(file, db)
			models.FillCatseyesBucket(db)
		case regexmap["ephemeris"].MatchString(filepath.Base(file)):
			tlemap := models.GetTLES(file)
			models.GetBeamplan(tlemap, bpfilelist, db)
			models.FleetTicker(db)
		default:
			continue
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
