package handlers

import (
	"github.com/alexmspina/worldmap/server/models"
	"github.com/functionalfoundry/graphqlws"
)

// GraphqlwsHandler web socket handler for graphql subscriptions
var GraphqlwsHandler = graphqlws.NewHandler(graphqlws.HandlerConfig{
	SubscriptionManager: models.SubscriptionManager,

	Authenticate: func(authToken string) (interface{}, error) {
		return "Fleet", nil
	},
})
