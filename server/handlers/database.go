package handlers

import (
	"fmt"
	"net/http"

	"github.com/alexmspina/worldmap/server/models"
	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
)

// DBHandler custom handler to accept bolt database
type DBHandler struct {
	Handler func(db *bolt.DB, w http.ResponseWriter, r *http.Request)
	DB      *bolt.DB
}

func (dbhandler *DBHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dbhandler.Handler(dbhandler.DB, w, r)
}

// GetStuff2 another getstuff function for testing passing db to handler
func GetStuff2(db *bolt.DB) httprouter.Handle {
	sats := []string{"M001", "M002", "M003", "M004", "M005", "M006", "M007", "M008", "M009", "M010", "M011", "M012", "M013", "M014", "M015", "M016"}
	bytes := make([]byte, 0)

	for _, s := range sats {
		sat := models.GetDBObject(s, db, "DB", "SATPOS")
		AppendBytes(&bytes, sat)
	}

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Fprintf(w, string(bytes))
	}
}
