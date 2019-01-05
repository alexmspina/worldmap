package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"

	"github.com/alexmspina/worldmap/server/helpers"
	"github.com/boltdb/bolt"
	"github.com/graphql-go/graphql"
)

// CatseyeFeature struct modeling geojson polygon struct
type CatseyeFeature struct {
	Type       string          `json:"type"`
	Geometry   PolygonGeometry `json:"geometry"`
	Properties ZoneProperties  `json:"properties"`
}

// CatseyeType graphql object for catseye features
var CatseyeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Catseye",
	Fields: graphql.Fields{
		"type": &graphql.Field{
			Type: graphql.String,
		},
		"geometry": &graphql.Field{
			Type:        PolyGeoType,
			Description: "coordinates that build catseye",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				s := params.Source.(CatseyeFeature)

				return s.Geometry, nil
			},
		},
		"properties": &graphql.Field{
			Type:        ZonePropsType,
			Description: "zone properties",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				s := params.Source.(CatseyeFeature)

				return s.Properties, nil
			},
		},
	},
})

// CatseyeFeatureCollection struct modeling catseye feature geojson
type CatseyeFeatureCollection struct {
	Type     string           `json:"type"`
	Features []CatseyeFeature `json:"features"`
}

// CatseyeFeatureCollectionType graphql object for catseye feature collections
var CatseyeFeatureCollectionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CatseyeFeatureCollection",
	Fields: graphql.Fields{
		"type": &graphql.Field{
			Type: graphql.String,
		},
		"features": &graphql.Field{
			Type:        graphql.NewList(CatseyeType),
			Description: "catseye features",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				s := params.Source.(CatseyeFeatureCollection)

				return s.Features, nil
			},
		},
	},
})

// ZoneFeature struct that models json object for zones
type ZoneFeature struct {
	Type       string         `json:"type"`
	Properties ZoneProperties `json:"properties"`
}

// ZoneProperties struct that models json object for properties held by zone features
type ZoneProperties struct {
	Subregion string  `json:"subregion"`
	ZoneID    string  `json:"zoneid"`
	StartLng  float64 `json:"startlng"`
	CenterLng float64 `json:"centerlng"`
	EndLng    float64 `json:"endlng"`
	Gateway   string  `json:"gateway"`
}

// ZonePropsType graphql type for target feature properties
var ZonePropsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CatseyeProps",
	Fields: graphql.Fields{
		"subregion": &graphql.Field{
			Type: graphql.String,
		},
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"startLng": &graphql.Field{
			Type: graphql.Float,
		},
		"centerLng": &graphql.Field{
			Type: graphql.Float,
		},
		"endLng": &graphql.Field{
			Type: graphql.Float,
		},
		"gateway": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// FillZonesBucket initializes a BoltDB bucket from TARGETS file
func FillZonesBucket(f string) error {
	// Open file from filename
	r := OpenCSV(f)

	// Get column header values
	header := getHeader(r)
	var err error

	// Create list of feature structs
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if record[0] == header[0] {
			continue
		}

		f := buildZoneFeature(record)
		id := f.Properties.ZoneID

		featureBytes, err := json.MarshalIndent(f, "", "\t")
		if err != nil {
			panic(err)
		}
		featureBytes = bytes.Replace(featureBytes, []byte("\\u0026"), []byte("&"), -1)
		featureBytes = bytes.Trim(featureBytes, "\r")

		err = DB.Update(func(tx *bolt.Tx) error {
			err = tx.Bucket([]byte("DB")).Bucket([]byte("ZONES")).Put([]byte(id), featureBytes)
			if err != nil {
				return fmt.Errorf("could not fill zones bucket: %v", err)
			}
			return nil
		})
	}
	fmt.Println("Zones bucket filled")

	return err
}

// FillCatseyesBucket fills bolt db with catseye polygon geojson objects
func FillCatseyesBucket() {
	catseyes := make([]CatseyeFeature, 0)
	// Get zones from db and calculate coordinates for catseye polygon
	DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB")).Bucket([]byte("ZONES"))
		b.ForEach(func(k, v []byte) error {
			var z ZoneFeature
			json.Unmarshal(v, &z)

			centerLng := z.Properties.CenterLng
			startLng := z.Properties.StartLng
			endLng := z.Properties.EndLng

			if endLng < centerLng {
				endLng = endLng + 360
			}

			if startLng > centerLng {
				startLng = startLng - 360
			}

			startPoint := []float64{0, startLng}
			centerPoint := []float64{0, centerLng}
			endPoint := []float64{0, endLng}

			// fmt.Println("start", startPoint)
			// fmt.Println("center", centerPoint)
			// fmt.Println("end", endPoint)

			catseyepoints := make([][]float64, 0)
			ComputeCoverageCircle(startPoint, centerPoint, "start", &catseyepoints)
			ComputeCoverageCircle(endPoint, centerPoint, "end", &catseyepoints)

			// lastPoint :=

			cat := BuildCatseyeFeature(catseyepoints, z.Properties)
			catseyes = append(catseyes, cat)
			return nil
		})
		return nil
	})

	// Put new catseye features in db
	for _, eye := range catseyes {
		id := eye.Properties.ZoneID

		featureBytes, err := json.MarshalIndent(eye, "", "\t")
		if err != nil {
			panic(err)
		}
		featureBytes = bytes.Replace(featureBytes, []byte("\\u0026"), []byte("&"), -1)
		featureBytes = bytes.Trim(featureBytes, "\r")

		err = DB.Update(func(tx *bolt.Tx) error {
			err = tx.Bucket([]byte("DB")).Bucket([]byte("CATSEYES")).Put([]byte(id), featureBytes)
			if err != nil {
				return fmt.Errorf("could not fill catseyes bucket: %v", err)
			}
			return nil
		})
	}
	fmt.Println("Catseyes bucket filled")
}

// BuildCatseyeFeature creates a catseye struct
func BuildCatseyeFeature(c [][]float64, p ZoneProperties) CatseyeFeature {
	g := PolygonGeometry{
		"Polygon",
		c,
	}
	v := CatseyeFeature{
		"Feature",
		g,
		p,
	}
	return v
}

func overLngWindow(initlng ...float64) []float64 {
	reslng := make([]float64, 0)
	for _, l := range initlng {
		if l > 180.0 {
			l = l - 360
			reslng = append(reslng, l)
		} else {
			reslng = append(reslng, l)
		}
	}
	return reslng
}

func buildZoneFeature(r []string) ZoneFeature {
	StartLng := helpers.ConvertStringToFloat64(r[2])
	CenterLng := helpers.ConvertStringToFloat64(r[3])
	EndLng := helpers.ConvertStringToFloat64(r[4])

	// NewLngs := overLngWindow(StartLng, CenterLng, EndLng)

	// StartLng = NewLngs[0]
	// CenterLng = NewLngs[1]
	// EndLng = NewLngs[2]

	props := ZoneProperties{
		Subregion: r[0],
		ZoneID:    r[1],
		StartLng:  StartLng,
		CenterLng: CenterLng,
		EndLng:    EndLng,
		Gateway:   r[5],
	}
	f := ZoneFeature{
		"Zone",
		props,
	}
	return f
}

// ComputeCoverageCircle generate list of lat/lng points
func ComputeCoverageCircle(p []float64, c []float64, s string, l *[][]float64) {
	elevation := helpers.Degs2Rads(10.0)
	height := 8062000.0
	earthRadius := 6378135.0
	subSatLat := helpers.Degs2Rads(p[0])
	subSatLng := helpers.Degs2Rads(p[1])
	centralAngle := math.Acos(math.Cos(elevation)/(1+(height/earthRadius))) - elevation
	centralAngleDeg := helpers.Rads2Degs(centralAngle)
	beamRadius := 2.0 * math.Pi * earthRadius * (centralAngleDeg / 360.0)

	for i := 0; i < 360; i++ {
		j := float64(i)
		j = helpers.Degs2Rads(j)
		latrad := math.Asin(math.Sin(subSatLat)*math.Cos(beamRadius/earthRadius) + math.Cos(subSatLat)*math.Sin(beamRadius/earthRadius)*math.Cos(j))
		lngrad := subSatLng + math.Atan2(math.Cos(beamRadius/earthRadius)-math.Sin(subSatLat)*math.Sin(latrad), math.Sin(j)*math.Sin(beamRadius/earthRadius)*math.Cos(subSatLat))
		latdeg := helpers.Rads2Degs(latrad)
		lngdeg := helpers.Rads2Degs(lngrad)
		point := []float64{
			latdeg,
			lngdeg,
		}

		switch s {
		case "full":
			*l = append(*l, point)
		case "start":

			if lngdeg > c[1] {
				// fmt.Println("center", c[1], "start lng", lngdeg)
				*l = append(*l, point)
			}
		case "end":
			if lngdeg < c[1] {
				// fmt.Println("center", c[1], "end lng", lngdeg)
				*l = append(*l, point)
			}
		default:
			continue
		}
	}
}

// GetCurrentZone determine which zone the satellite is currently servicing
func GetCurrentZone(satlng float64) []string {
	var zoneid []string
	var satlngadjusted float64
	err := DB.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB")).Bucket([]byte("ZONES"))
		b.ForEach(func(k, v []byte) error {
			var zone ZoneFeature
			json.Unmarshal(v, &zone)
			zonestartlng := zone.Properties.StartLng
			zoneendlng := zone.Properties.EndLng

			if satlng < 0.0 {
				satlngadjusted = satlng + 360.0
			} else {
				satlngadjusted = satlng
			}

			if zonestartlng > zoneendlng {
				zoneendlng = zoneendlng + 360.0

				if zoneendlng > 360.0 {
					satlngadjusted = satlng + 360
				}

				if satlngadjusted > zonestartlng && satlngadjusted < zoneendlng {
					zoneid = append(zoneid, string(k))
				}
			} else {
				if satlngadjusted > zonestartlng && satlngadjusted < zoneendlng {
					zoneid = append(zoneid, string(k))
				}
			}

			return nil
		})
		return nil
	})
	helpers.PanicErrors(err)
	return zoneid
}

// GetCatseye queries bolt db for the desired target
func GetCatseye(s string) CatseyeFeature {
	var catseyefeature CatseyeFeature
	err := DB.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB")).Bucket([]byte("CATSEYES"))
		catseye := b.Get([]byte(s))
		json.Unmarshal(catseye, &catseyefeature)

		return nil
	})
	helpers.PanicErrors(err)

	return catseyefeature
}

// GetAllCatseyes grabs all the catseyes from the CATSEYE bucket
func GetAllCatseyes() []CatseyeFeature {
	var catseyeFeatureList []CatseyeFeature
	err := DB.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB")).Bucket([]byte("CATSEYES"))
		b.ForEach(func(k, v []byte) error {
			var catseyeFeature CatseyeFeature
			json.Unmarshal(v, &catseyeFeature)
			catseyeFeatureList = append(catseyeFeatureList, catseyeFeature)
			return nil
		})
		return nil
	})
	helpers.PanicErrors(err)

	return catseyeFeatureList
}
