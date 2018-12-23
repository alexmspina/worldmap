package models

import (
	"fmt"

	"github.com/graphql-go/graphql"
)

// http://localhost:8080/satellite?query={satellite(id:%22M007%22){id,latitude,longitude,velocity,altitude,mission{id,config,gatewayID,gatewayOBAnt,gatewayMaxPointingTime,beams{id,epcs,targetOBAnt,targetMaxPointingTime,camp,campMode,campGain,ldla,ldlaMode,ldlaFcaGain,ldlaGcaGain,ldlaScaGain}}}}
// http://localhost:8080/target?query={target(id:%2235%22){geometry{coordinates},properties{shortName}}
// http://localhost:8080/catseye?query={catseye(id:%2210%22){geometry{coordinates},properties{subregion}}
// http://localhost:8080/targets?query={targets{geometry{coordinates},properties{shortName}}}

// RootQuery main graphql query for schema
var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
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

				return SatelliteFeature{}, nil
			},
		},
		"satelliteFeatureCollection": &graphql.Field{
			Type:        SatelliteFeatureCollectionType,
			Description: "Get all targets and their properties",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				fmt.Println("REQUESTING SATELLITES...")
				satellites := GetMovingSatellites()
				satelliteFeatureCollection := SatelliteFeatureCollection{
					Type:     "featureCollection",
					Features: satellites,
				}
				return satelliteFeatureCollection, nil
			},
		},
		"target": &graphql.Field{
			Type:        TargetType,
			Description: "Get a single target and its properties",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				idQuery, isOK := params.Args["id"].(string)
				if isOK {
					target := GetTarget(idQuery)
					return target, nil
				}

				return TargetFeature{}, nil
			},
		},
		"targetFeatureCollection": &graphql.Field{
			Type:        TargetFeatureCollectionType,
			Description: "Get all targets and their properties",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				fmt.Println("REQUESTING TARGETS...")
				targets := GetTargets()
				targetFeatureCollection := TargetFeatureCollection{
					Type:     "featureCollection",
					Features: targets,
				}
				return targetFeatureCollection, nil
			},
		},
		"targets": &graphql.Field{
			Type:        graphql.NewList(TargetType),
			Description: "Get all targets and their properties",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				fmt.Println("REQUESTING TARGETS...")
				targets := GetTargets()
				return targets, nil
			},
		},
		"catseye": &graphql.Field{
			Type:        CatseyeType,
			Description: "Get a single target and its properties",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				idQuery, isOK := params.Args["id"].(string)
				if isOK {
					target := GetCatseye(idQuery)
					return target, nil
				}

				return CatseyeFeature{}, nil
			},
		},
	},
})

// ExecuteQuery performs a graphql query
func ExecuteQuery(query string, params graphql.Params) *graphql.Result {
	result := graphql.Do(params)
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}
