package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/alexmspina/worldmap/server/handlers"
	"github.com/alexmspina/worldmap/server/helpers"
	"github.com/alexmspina/worldmap/server/models"
	"github.com/julienschmidt/httprouter"
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

	// Setup bolt database
	db, _ := models.SetupDB()

	// Process the selected files depending on their type and fill bolt db buckets
	models.ProcessInitFiles(files, regexmap, db)

	// Process files if they are tles
	sgp4sats := models.ProcessEphemeris(files, regexmap, db, bpfilelist)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go models.FleetTicker(ticker, sgp4sats, db)

	router := httprouter.New()
	router.GET("/", handlers.Index)
	router.GET("/payloadmissions", handlers.GetStuff2(db))

	log.Fatal(http.ListenAndServe(":8080", router))
}
