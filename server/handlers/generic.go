package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Index initial path resolved by server
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "C:/Users/aspina/go/src/github.com/alexmspina/worldmap/server/build/index.html")
}
