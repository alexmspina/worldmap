package main

import (
	"encoding/json"
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
	"github.com/graphql-go/graphql"
	"github.com/julienschmidt/httprouter"
)

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"satellite": &graphql.Field{
			Type:        models.SatelliteType,
			Description: "Get a single satellite, its location, and current mission",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				idQuery, isOK := params.Args["id"].(string)
				if isOK {
					livesat := models.GetSatellitePosition(idQuery)
					return livesat, nil
				}

				return models.SatelliteInMotion{}, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
})

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

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
	// db, _ := models.SetupDB()

	// Process the selected files depending on their type and fill bolt db buckets
	models.ProcessInitFiles(files, regexmap)

	// Process files if they are tles
	sgp4sats := models.ProcessEphemeris(files, regexmap, bpfilelist)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go models.FleetTicker(ticker, sgp4sats)

	router := httprouter.New()
	router.GET("/", handlers.Index)
	router.GET("/graphql", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})
	// router.GET("/payloadmissions", handlers.GetStuff2(db))

	log.Fatal(http.ListenAndServe(":8080", router))
}
