package models

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"

	satellite "github.com/joshuaferrara/go-satellite"
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
func ProcessInitFiles(files []string, regexmap map[string]*regexp.Regexp) {
	for _, file := range files {
		switch true {
		case regexmap["TARGETS"].MatchString(filepath.Base(file)):
			FillTargetsBucket(file)
		case regexmap["ZONES"].MatchString(filepath.Base(file)):
			FillZonesBucket(file)
			FillCatseyesBucket()
		default:
			continue
		}
	}
}

// ProcessEphemeris takes tle and beamplan and initializes db entries
func ProcessEphemeris(files []string, regexmap map[string]*regexp.Regexp, bpfilelist map[string]string) map[string]satellite.Satellite {
	sgp4sats := make(map[string]satellite.Satellite, 0)
	for _, file := range files {
		switch true {
		case regexmap["ephemeris"].MatchString(filepath.Base(file)):
			tlemap := GetTLES(file)
			GetBeamplan(tlemap, bpfilelist)
			satStates := GetSatelliteStates()
			sgp4sats = InitSatelliteSGP4(satStates)
		default:
			continue
		}
	}
	return sgp4sats
}
