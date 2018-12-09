package models

// PolygonGeometry struct that models geojson geometry type for polygons
type PolygonGeometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

// PointGeometry struct that models geojson geometry type for points
type PointGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}
