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
)

// CatseyeFeature struct modeling geojson polygon struct
type CatseyeFeature struct {
	Type       string          `json:"type"`
	Geometry   PolygonGeometry `json:"geometry"`
	Properties ZoneProperties  `json:"properties"`
}

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
			centerPoint := []float64{0, z.Properties.CenterLng}
			endPoint := []float64{0, z.Properties.StartLng}
			startPoint := []float64{0, z.Properties.EndLng}

			catseyepoints := make([][]float64, 0)
			ComputeCoverageCircle(endPoint, centerPoint, "end", &catseyepoints)
			ComputeCoverageCircle(startPoint, centerPoint, "start", &catseyepoints)

			cat := BuildCatseyeFeature(catseyepoints, z.Properties)
			catseyes = append(catseyes, cat)
			return nil
		})
		return nil
	})

	// Put new catseye features in db
	for _, eye := range catseyes {
		id := eye.Properties.Subregion

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
		"Catseye",
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

	NewLngs := overLngWindow(StartLng, CenterLng, EndLng)

	StartLng = NewLngs[0]
	CenterLng = NewLngs[1]
	EndLng = NewLngs[2]

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
	elevation := helpers.Degs2Rads(0)
	height := 8062000.0
	earthRadius := 6378000.0
	subSatLat := helpers.Degs2Rads(p[0])
	subSatLng := helpers.Degs2Rads(p[1])
	centralAngle := math.Acos(math.Cos(elevation)/(1+height/earthRadius)) - elevation

	for i := 0; i < 360; i++ {
		j := float64(i)
		lat := helpers.Rads2Degs(math.Asin(math.Sin(subSatLat)*math.Cos(centralAngle) + math.Cos(subSatLat)*math.Sin(centralAngle)*math.Cos(j)))
		lng := helpers.Rads2Degs(subSatLng + math.Atan2(math.Sin(j)*math.Sin(centralAngle)*math.Cos(subSatLat), math.Cos(centralAngle)-math.Sin(subSatLat)*math.Sin(math.Asin(math.Sin(subSatLat)*math.Cos(centralAngle)+math.Cos(subSatLat)*math.Sin(centralAngle)*math.Cos(j)))))
		point := []float64{
			lat,
			lng,
		}

		switch s {
		case "full":
			*l = append(*l, point)
		case "start":
			if lng < c[1] {
				*l = append(*l, point)
			}
		case "end":
			if lng > c[1] {
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
	err := DB.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB")).Bucket([]byte("ZONES"))
		b.ForEach(func(k, v []byte) error {
			var zone ZoneFeature
			json.Unmarshal(v, &zone)
			zonestartlng := zone.Properties.StartLng
			zoneendlng := zone.Properties.EndLng

			// shift longitudes less than 0 to 0 - 360 range for easy zone placement
			if zoneendlng < zonestartlng {
				zoneendlng = zoneendlng + 360.0

				if satlng < 0 {
					satlngadjusted := satlng + 360.0

					if satlngadjusted > zonestartlng && satlngadjusted < zoneendlng {
						zoneid = append(zoneid, string(k))
					}
				}
			} else {
				if satlng > zonestartlng && satlng < zoneendlng {
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
