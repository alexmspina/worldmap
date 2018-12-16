package models

import (
	"github.com/graphql-go/graphql"
)

// PolygonGeometry struct that models geojson geometry type for polygons
type PolygonGeometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

// PolyGeoType graphql object for individual catseye queries
var PolyGeoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "polygonGeometry",
	Fields: graphql.Fields{
		"type": &graphql.Field{
			Type: graphql.String,
		},
		"coordinates": &graphql.Field{
			Type:        graphql.NewList(graphql.NewList(graphql.Float)),
			Description: "List of coordinates to build shape",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				s := params.Source.(PolygonGeometry)

				return s.Coordinates, nil
			},
		},
	},
})

// PointGeometry struct that models geojson geometry type for points
type PointGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// PointGeoType graphql object for individual beamplan mission queries
var PointGeoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "polygonGeometry",
	Fields: graphql.Fields{
		"type": &graphql.Field{
			Type: graphql.String,
		},
		"coordinates": &graphql.Field{
			Type:        graphql.NewList(graphql.Float),
			Description: "List of coordinates to build shape",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				s := params.Source.(PointGeometry)

				return s.Coordinates, nil
			},
		},
	},
})
