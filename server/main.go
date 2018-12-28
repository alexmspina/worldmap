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

	// add index.html to bld directory
	index := helpers.Join(*bld, "/index.html")
	fmt.Println(index)
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
