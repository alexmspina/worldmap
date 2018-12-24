package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/alexmspina/worldmap/server/appmount"
	"github.com/alexmspina/worldmap/server/handlers"
	"github.com/julienschmidt/httprouter"
)

func main() {
	// parse command-line flag to determine root directory location of necessary files
	dir := flag.String("dir", "No data directory provided", "input the directory where the initial data files are located")
	// bld := flag.String("bld", "No duild directory provided", "input the directory where the build files are located")

	flag.Parse()

	tickerChannel := time.NewTicker(time.Second).C
	go appmount.AppMount(tickerChannel, dir)

	graphqlHandler := http.HandlerFunc(handlers.GraphqlHandlerFunc)
	router := httprouter.New()

	router.GET("/", handlers.Index)
	router.POST("/graphql", handlers.DisableCors(graphqlHandler))
	// router.ServeFiles(*bld, http.Dir("index.html"))

	router.ServeFiles("/build/*filepath", http.Dir("static"))

	log.Fatal(http.ListenAndServe(":8080", router))
	// log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir(*bld))))
}
