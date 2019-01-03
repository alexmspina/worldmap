package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexmspina/worldmap/server/appmount"
	"github.com/alexmspina/worldmap/server/handlers"
	"github.com/alexmspina/worldmap/server/helpers"
	"github.com/alexmspina/worldmap/server/models"
	"github.com/julienschmidt/httprouter"
)

func main() {
	// parse command-line flag to determine root directory location of necessary files
	dir := flag.String("dir", "No data directory provided", "input the directory where the initial data files are located")
	bld := flag.String("bld", "No duild directory provided", "input the directory where the build files are located")
	flag.Parse()

	// mount app
	tickerChannel := time.NewTicker(time.Second).C
	go appmount.AppMount(tickerChannel, dir)

	ro4 := models.GetCatseye("22")
	lowlat := ro4.Geometry.Coordinates[0][0]
	hilat := ro4.Geometry.Coordinates[0][0]
	lowlng := ro4.Geometry.Coordinates[0][1]
	hilng := ro4.Geometry.Coordinates[0][1]
	for _, i := range ro4.Geometry.Coordinates {
		if i[0] > hilat {
			hilat = i[0]
		}
		if i[0] < lowlat {
			lowlat = i[0]
		}
		if i[1] > hilng {
			hilng = i[1]
		}
		if i[1] < lowlng {
			lowlng = i[1]
		}
	}

	fmt.Printf("hilat: %v\n", hilat)
	fmt.Printf("lowlat: %v\n", lowlat)
	fmt.Printf("hilng: %v\n", hilng)
	fmt.Printf("lowlng: %v\n", lowlng)

	// add index.html to bld directory
	index := helpers.Join(*bld, "/index.html")
	static := helpers.Join(*bld, "/static")

	// http router with
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.ServeFile(w, r, index)
	})
	graphqlHandler := http.HandlerFunc(handlers.GraphqlHandlerFunc)
	router.POST("/graphql", handlers.DisableCors(graphqlHandler))
	router.ServeFiles("/static/*filepath", http.Dir(static))
	log.Fatal(http.ListenAndServe(":8080", router))
}
