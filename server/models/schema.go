package models

import "github.com/graphql-go/graphql"

// Schema graphql schema
var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: RootQuery,
})
