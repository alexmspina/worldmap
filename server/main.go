package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/alexmspina/worldmap/server/appmount"
	"github.com/alexmspina/worldmap/server/handlers"
	"github.com/alexmspina/worldmap/server/models"
	"github.com/julienschmidt/httprouter"
)

func main() {
	// parse command-line flag to determine root directory location of necessary files
	dir := flag.String("dir", "No directory provided", "input the directory where the initial files are located")
	appmount.ParseFlag(dir)

	tickerChannel := time.NewTicker(time.Second).C
	go appmount.AppMount(tickerChannel, dir)

	router := httprouter.New()
	router.GET("/", handlers.Index)
	router.GET("/subscriptions", handlers.WrapHandler(handlers.GraphqlwsHandler))
	router.GET("/satellite", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		result := models.ExecuteQuery(r.URL.Query().Get("query"), models.Schema)
		json.NewEncoder(w).Encode(result)
	})
	router.GET("/satellites", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		result := models.ExecuteQuery(r.URL.Query().Get("query"), models.Schema)
		json.NewEncoder(w).Encode(result)
	})
	router.GET("/target", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		result := models.ExecuteQuery(r.URL.Query().Get("query"), models.Schema)
		json.NewEncoder(w).Encode(result)
	})
	router.GET("/targets", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		result := models.ExecuteQuery(r.URL.Query().Get("query"), models.Schema)
		json.NewEncoder(w).Encode(result)
	})
	router.GET("/catseye", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		result := models.ExecuteQuery(r.URL.Query().Get("query"), models.Schema)
		json.NewEncoder(w).Encode(result)
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
