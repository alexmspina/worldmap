package appmount

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/alexmspina/worldmap/server/helpers"
	"github.com/alexmspina/worldmap/server/models"
	satellite "github.com/joshuaferrara/go-satellite"
)

// AppMount initializes app state
func AppMount(t <-chan time.Time, dir *string) {
	fmt.Println("Mounting application")

	// use Walk function to traverse root directory provided and create a list of files
	files := make([]string, 0)
	helpers.GetFilesFromDirectory(&files, *dir)

	// Create map of regular expressions
	regexmap := make(map[string]*regexp.Regexp, 0)
	preregexlist := []string{"TARGETS", "BEAMPLAN_LONGFORMAT", "ZONES", "ephemeris"}
	helpers.CreateRegexp(regexmap, preregexlist)

	bpregexmap := make(map[string]*regexp.Regexp, 0)
	bppreregex := []string{"mute", "B3", "M001", "M013"}
	helpers.CreateRegexp(bpregexmap, bppreregex)

	bpfilelist := make(map[string]string, 0)
	models.GetBeamplanFiles(files, regexmap["BEAMPLAN_LONGFORMAT"], bpregexmap, bpfilelist)

	fmt.Println("Building data models")

	// Process the selected files depending on their type and fill bolt db buckets
	models.ProcessInitFiles(files, regexmap)

	// Process files if they are tles
	sgp4sats := models.ProcessEphemeris(files, regexmap, bpfilelist)

	for {
		select {
		case <-t:
			currentTime := <-t
			models.UpdateSatPos(currentTime, sgp4sats)
		}
	}
}

// AppTicker global ticker for entire app
func AppTicker(ticker <-chan *time.Ticker, sgp4sats map[string]satellite.Satellite) {

}

// ParseFlag takes a flag pointer to a string and parses it
func ParseFlag(dir *string) {
	flag.Parse()
	if *dir == "No directory provided" {
		fmt.Println("No directory provided for initial setup.")
		os.Exit(1)
	}

	fmt.Println("Thank you for providing a directory.")
}
