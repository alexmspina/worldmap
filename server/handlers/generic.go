package handlers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Index initial path resolved by server
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	fmt.Fprint(w, "Welcome!\n")
}
