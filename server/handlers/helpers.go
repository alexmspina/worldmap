package handlers

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// DisableCors disables cors filtered requests and allows options requests from graphql client
func DisableCors(h http.Handler) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		// result := models.ExecuteQuery(r.URL.Query().Get("query"), models.Schema)
		// json.NewEncoder(w).Encode(result)
		h.ServeHTTP(w, r)
	})
}

// StringKey creates a special type so it doesn't conflict with standard strings
type StringKey string

// WrapHandler wraps a standard http handler so it can be used with httprouter package
func WrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Take the context out from the request
		ctx := r.Context()

		var paramskey StringKey
		paramskey = "params"

		// Get new context with key-value "params" -> "httprouter.Params"
		ctx = context.WithValue(ctx, paramskey, ps)

		// Get new http.Request with the new context
		r = r.WithContext(ctx)

		// Call your original http.Handler
		h.ServeHTTP(w, r)
	}
}
