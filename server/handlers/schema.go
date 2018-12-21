package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/alexmspina/worldmap/server/models"
)

// H test graphql query handler
var H = handler.New(&handler.Config{
	Schema: &models.Schema,
	Pretty: true,
})

// GraphqlHandlerFunc handler func that uses graphql-go handler
func GraphqlHandlerFunc(w http.ResponseWriter, r *http.Request) {
	// get query
	opts := handler.NewRequestOptions(r)

	// execute graphql query
	params := graphql.Params{
		Schema:         models.Schema, // defined in another file
		RequestString:  opts.Query,
		VariableValues: opts.Variables,
		OperationName:  opts.OperationName,
		Context:        r.Context(),
	}

	result := models.ExecuteQuery(r.URL.Query().Get("query"), params)

	json.NewEncoder(w).Encode(result)
}
