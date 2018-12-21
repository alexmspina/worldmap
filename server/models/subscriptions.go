package models

import (
	"fmt"

	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql"
)

// SubscriptionManager create a subscription manager
var SubscriptionManager = graphqlws.NewSubscriptionManager(&Schema)

// RootSubscription main graphql query for schema
var RootSubscription = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootSubscription",
	Fields: graphql.Fields{
		"satellite": &graphql.Field{
			Type:        SatelliteType,
			Description: "Get a single satellite, its location, and current mission",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				idQuery, isOK := params.Args["id"].(string)
				if isOK {
					livesat := GetSatellitePosition(idQuery)
					return livesat, nil
				}

				return SatelliteInMotion{}, nil
			},
		},
	},
})

// ProcessSubscriptions rotate through subscriptions to check it anything updated and should be sent to client
func ProcessSubscriptions() {
	// This assumes you have access to the above subscription manager
	subscriptions := SubscriptionManager.Subscriptions()

	for _, conn := range subscriptions {
		// Things you have access to here:
		fmt.Println("hello", conn)

		// for _, subscription := range subscriptions[conn] {
		// 	// Things you have access to here:
		// 	subscription.ID            // The subscription ID (unique per conn)
		// 	subscription.OperationName // The name of the subcription
		// 	subscription.Query         // The subscription query/queries string
		// 	subscription.Variables     // The subscription variables
		// 	subscription.Document      // The GraphQL AST for the subscription
		// 	subscription.Fields        // The names of top-level queries
		// 	subscription.Connection    // The GraphQL WS connection

		// 	// Prepare an execution context for running the query
		// 	ctx := context.Context()

		// 	// Re-execute the subscription query
		// 	params := graphql.Params{
		// 		Schema:         Schema, // The GraphQL schema
		// 		RequestString:  subscription.Query,
		// 		VariableValues: subscription.Variables,
		// 		OperationName:  subscription.OperationName,
		// 		Context:        ctx,
		// 	}
		// 	result := graphql.Do(params)

		// 	// Send query results back to the subscriber at any point
		// 	data := graphqlws.DataMessagePayload{
		// 		// Data can be anything (interface{})
		// 		Data: result.Data,
		// 		// Errors is optional ([]error)
		// 		Errors: graphqlws.ErrorsFromGraphQLErrors(result.Errors),
		// 	}
		// 	subscription.SendData(&data)
		// }
	}
}
