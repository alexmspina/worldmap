package handlers

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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
