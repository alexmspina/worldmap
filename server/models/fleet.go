package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/alexmspina/worldmap/server/helpers"
	"github.com/graphql-go/graphql"

	"github.com/boltdb/bolt"
	satellite "github.com/joshuaferrara/go-satellite"
)

// Fleet struct of all satellite states derived from the tle and beamplan
type Fleet struct {
	Satellites map[string]SatelliteState `json:"satellites"`
}

// SatelliteState struct that models the static info of a satellite including tle lines and current beamplan
type SatelliteState struct {
	TLELine1 string            `json:"tleLine1"`
	TLELine2 string            `json:"tleLine2"`
	Missions []BeamplanMission `json:"missions"`
}

// BeamplanMission struct modeling beamplan mission
type BeamplanMission struct {
	ID                     string           `json:"id"`
	MissionConfig          string           `json:"missionConfig"`
	GatewayTargetID        string           `json:"gatewayTargetID"`
	GatewayOBAntID         string           `json:"gatewayOBAntID"`
	GatewayPointingMaxTime string           `json:"gatewayPointingMaxTime"`
	Beams                  []BeamProperties `json:"beams"`
}

// MissionType graphql object for individual beamplan mission queries
var MissionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mission",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"config": &graphql.Field{
			Type: graphql.String,
		},
		"gatewayID": &graphql.Field{
			Type: graphql.String,
		},
		"gatewayOBAnt": &graphql.Field{
			Type: graphql.String,
		},
		"gatewayMaxPointingTime": &graphql.Field{
			Type: graphql.String,
		},
		"beams": &graphql.Field{
			Type:        graphql.NewList(BeamPropsType),
			Description: "Get the beams from the current mission",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				s := params.Source.(BeamplanMission)

				return s.Beams, nil
			},
		},
	},
})

// BeamProperties struct modeling individual beam settings
type BeamProperties struct {
	ID                    string `json:"id"`
	EPCList               string `json:"epcList"`
	TargetOBAntID         string `json:"targetOBAntID"`
	TargetMaxPointingTime string `json:"targetMaxPointingTime"`
	CampID                string `json:"campID"`
	CampMode              string `json:"campMode"`
	CampGain              string `json:"campGain"`
	LDLAID                string `json:"ldlaID"`
	LDLAMode              string `json:"ldlaMode"`
	LDLAFCAGain           string `json:"ldlaFCAGain"`
	LDLAGCAGain           string `json:"ldlaGCAGain"`
	LDLASCAGain           string `json:"ldlaSCAGain"`
}

// BeamPropsType graphql object for Beam property queries
var BeamPropsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "BeamProperties",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"epcs": &graphql.Field{
			Type: graphql.String,
		},
		"targetOBAnt": &graphql.Field{
			Type: graphql.String,
		},
		"targetMaxPointingTime": &graphql.Field{
			Type: graphql.String,
		},
		"camp": &graphql.Field{
			Type: graphql.String,
		},
		"campMode": &graphql.Field{
			Type: graphql.String,
		},
		"campGain": &graphql.Field{
			Type: graphql.String,
		},
		"ldla": &graphql.Field{
			Type: graphql.String,
		},
		"ldlaMode": &graphql.Field{
			Type: graphql.String,
		},
		"ldlaFcaGain": &graphql.Field{
			Type: graphql.String,
		},
		"ldlaGcaGain": &graphql.Field{
			Type: graphql.String,
		},
		"ldlaScaGain": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// SatelliteInMotion structure to hold location and velocity data that is updated in real time and pushed to db
type SatelliteInMotion struct {
	ID        string            `json:"satelliteID"`
	Latitude  float64           `json:"latitude"`
	Longitude float64           `json:"longitude"`
	Velocity  float64           `json:"velocity"`
	Altitude  float64           `json:"altitude"`
	Mission   []BeamplanMission `json:"mission"`
}

// SatelliteType graphql object for satellite queries
var SatelliteType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Satellite",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"latitude": &graphql.Field{
			Type: graphql.Float,
		},
		"longitude": &graphql.Field{
			Type: graphql.Float,
		},
		"velocity": &graphql.Field{
			Type: graphql.Float,
		},
		"altitude": &graphql.Field{
			Type: graphql.Float,
		},
		"mission": &graphql.Field{
			Type:        graphql.NewList(MissionType),
			Description: "Get the missions currently serviced by the satellite",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				s := params.Source.(SatelliteInMotion)
				currentZones := GetCurrentZone(s.Longitude)
				currentMissions := GetCurrentMission(s.ID, currentZones)
				fmt.Println(currentMissions, "\n", len(currentMissions))
				return currentMissions, nil
			},
		},
	},
})

// SetCurrentMission receives a pointer to a satelliteInMotion and sets its current mission
func (s *SatelliteInMotion) SetCurrentMission(m []BeamplanMission) {
	s.Mission = m
}

// GetTLES creats a map of tles with map of tle lines
func GetTLES(tle string) map[string]map[string]string {
	// get pointer to file of tle
	tlefile, err := os.Open(tle)
	helpers.PanicErrors(err)

	// read file into byte slices
	tlereader := io.Reader(tlefile)
	tlelines, err := ioutil.ReadAll(tlereader)
	helpers.PanicErrors(err)

	// convert tle byte slice to string and clean it up
	tlestringlines := string(tlelines)
	cleantlelines := strings.Split(tlestringlines, "\n")

	// Determine if the tle has a title or starts at first satellite
	if len(cleantlelines[len(cleantlelines)-1]) <= 0 {
		cleantlelines = cleantlelines[:len(cleantlelines)-1]
	}
	if len(cleantlelines)%3 != 0 {
		cleantlelines = cleantlelines[1:]
	}

	// create map of tles sorted by satellite name
	tlemap := make(map[string]map[string]string, 0)
	for i := 0; i < len(cleantlelines)/3; i++ {
		tmpmap := make(map[string]string)
		tmpmap["firstline"] = strings.Trim(cleantlelines[i*3+1], "\r")
		tmpmap["secondline"] = strings.Trim(cleantlelines[i*3+2], "\r")
		o3bname := strings.Trim(cleantlelines[i*3], "\r")
		o3bnamelen := len(o3bname)
		name := o3bname[o3bnamelen-4:]
		tlemap[name] = tmpmap
	}

	return tlemap
}

// GetBeamplan determines what beamplan file to use based on the satellite being processed
func GetBeamplan(tlemap map[string]map[string]string, bpfiles map[string]string) {
	spares := []string{"M002", "M004", "M005"}
	activenotb3 := []string{"M001", "M003", "M006", "M007", "M008", "M009", "M010", "M011", "M012"}
	activeb3 := []string{"M013", "M014", "M015", "M016"}

	for sat, tle := range tlemap {
		switch true {
		case helpers.StringInSlice(sat, spares):
			bpfile := bpfiles["SPARE"]
			satstate := BuildSatelliteState(bpfile, tle, sat)
			FillFleetBucket(sat, satstate)
		case helpers.StringInSlice(sat, activenotb3):
			switch sat {
			case "M001":
				bpfile := bpfiles["M001"]
				satstate := BuildSatelliteState(bpfile, tle, sat)
				FillFleetBucket(sat, satstate)
			default:
				bpfile := bpfiles["ACTIVE"]
				satstate := BuildSatelliteState(bpfile, tle, sat)
				FillFleetBucket(sat, satstate)
			}
		case helpers.StringInSlice(sat, activeb3):
			switch sat {
			case "M013":
				bpfile := bpfiles["M013"]
				satstate := BuildSatelliteState(bpfile, tle, sat)
				FillFleetBucket(sat, satstate)
			default:
				bpfile := bpfiles["B3"]
				satstate := BuildSatelliteState(bpfile, tle, sat)
				FillFleetBucket(sat, satstate)
			}
		}
	}
	fmt.Println("Fleet bucket filled.")
}

// FillFleetBucket initializes satellites from tle
func FillFleetBucket(satid string, satstate SatelliteState) error {
	satstatebytes, err := json.MarshalIndent(satstate, "", "\t")
	if err != nil {
		panic(err)
	}
	satstatebytes = bytes.Replace(satstatebytes, []byte("\\u0026"), []byte("&"), -1)
	satstatebytes = bytes.Trim(satstatebytes, "\r")

	err = DB.Update(func(tx *bolt.Tx) error {
		err = tx.Bucket([]byte("DB")).Bucket([]byte("FLEET")).Put([]byte(satid), satstatebytes)
		if err != nil {
			return fmt.Errorf("could not fill zones bucket: %v", err)
		}
		return nil
	})
	return nil
}

// BuildSatelliteState creates Fleet json struct
func BuildSatelliteState(bpfile string, tle map[string]string, satname string) SatelliteState {
	// open the beamplan file and pass to a reader
	bpreader := OpenCSV(bpfile)

	// Get column header values
	header := getHeader(bpreader)

	msns := make(map[string][][]string, 0)

	for {
		record, err := bpreader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			helpers.PanicErrors(err)
		}

		satid := record[0]
		msnid := record[1]
		if satid == header[0] {
			continue
		} else {
			switch satname {
			case satid:
				if msn, ok := msns[msnid]; ok {
					tmp := append(msn, record)
					msns[msnid] = tmp
				} else {
					tmp := make([][]string, 0)
					tmp = append(tmp, record)
					msns[msnid] = tmp
				}
			}
		}
	}

	msnsmap := make([]BeamplanMission, 0)
	for id, beams := range msns {
		tgts := make([]BeamProperties, 0)
		var msnconfig, gwtgtid, gwobantid, gwpntmxtime string
		for _, record := range beams {
			msnconfig = record[2]
			epclist := record[3]
			gwtgtid = record[5]
			gwobantid = record[6]
			gwpntmxtime = record[7]
			tgtid := record[8]
			tgtobantid := record[9]
			tgtpntmxtime := record[10]
			campid := record[11]
			campmode := record[12]
			campgain := record[13]
			ldlaid := record[14]
			ldlamode := record[15]
			ldlafcagain := record[16]
			ldlagcagain := record[17]
			ldlascagain := record[18]

			bmprops := BeamProperties{
				ID:                    tgtid,
				EPCList:               epclist,
				TargetOBAntID:         tgtobantid,
				TargetMaxPointingTime: tgtpntmxtime,
				CampID:                campid,
				CampMode:              campmode,
				CampGain:              campgain,
				LDLAID:                ldlaid,
				LDLAMode:              ldlamode,
				LDLAFCAGain:           ldlafcagain,
				LDLAGCAGain:           ldlagcagain,
				LDLASCAGain:           ldlascagain,
			}

			tgts = append(tgts, bmprops)
		}

		bpmsn := BeamplanMission{
			ID:                     id,
			MissionConfig:          msnconfig,
			GatewayTargetID:        gwtgtid,
			GatewayOBAntID:         gwobantid,
			GatewayPointingMaxTime: gwpntmxtime,
			Beams:                  tgts,
		}

		msnsmap = append(msnsmap, bpmsn)
	}

	satstate := SatelliteState{
		TLELine1: tle["firstline"],
		TLELine2: tle["secondline"],
		Missions: msnsmap,
	}

	return satstate
}

// GetSatellitePosition gets the current satellite position from SATPOS db
func GetSatellitePosition(s string) SatelliteInMotion {
	var livesat SatelliteInMotion
	err := DB.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB")).Bucket([]byte("SATPOS"))
		sat := b.Get([]byte(s))
		json.Unmarshal(sat, &livesat)

		return nil
	})
	helpers.PanicErrors(err)

	return livesat
}

// UpdateSatPos updates satellite positions using go rout
func UpdateSatPos(t time.Time, sgp4sats map[string]satellite.Satellite) {
	for i, sat := range sgp4sats {
		PropagateSatellite(t, sat, i)
	}
}

// PropagateSatellite take a satellite.Satellite struct and propagates it. Then pushes it to SatelliteInMotion channel
func PropagateSatellite(t time.Time, sat satellite.Satellite, id string) {
	utc := t.UTC()
	y, m, d := utc.Date()
	h, min, sec := utc.Clock()
	gmst := satellite.GSTimeFromDate(y, int(m), d, h, min, sec)
	pos, _ := satellite.Propagate(sat, y, int(m), d, h, min, sec)
	alt, vel, latlng := satellite.ECIToLLA(pos, gmst)
	latlngdeg := satellite.LatLongDeg(latlng)

	fmt.Println(id, latlngdeg.Longitude)

	satinmotion := SatelliteInMotion{
		ID:        id,
		Latitude:  latlngdeg.Latitude,
		Longitude: latlngdeg.Longitude,
		Velocity:  vel,
		Altitude:  alt,
	}

	FillSatPosBucket(satinmotion, id)
}

// GetCurrentMission gets the current mission from sat state object
func GetCurrentMission(satid string, missionids []string) []BeamplanMission {
	missions := make([]BeamplanMission, 0)

	err := DB.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB")).Bucket([]byte("FLEET"))
		sat := b.Get([]byte(satid))
		var satstate SatelliteState
		json.Unmarshal(sat, &satstate)

		for m := range missionids {
			missions = append(missions, satstate.Missions[m])
		}

		return nil
	})
	helpers.PanicErrors(err)

	return missions
}

// FillSatPosBucket fills the satellite position bucket with a satellite in motion object
func FillSatPosBucket(s SatelliteInMotion, id string) {

	satposBytes, err := json.MarshalIndent(s, "", "\t")
	helpers.PanicErrors(err)
	satposBytes = bytes.Replace(satposBytes, []byte("\\u0026"), []byte("&"), -1)
	satposBytes = bytes.Trim(satposBytes, "\r")

	err = DB.Batch(func(tx *bolt.Tx) error {
		err = tx.Bucket([]byte("DB")).Bucket([]byte("SATPOS")).Put([]byte(id), satposBytes)
		if err != nil {
			return fmt.Errorf("could not fill catseyes bucket: %v", err)
		}
		return nil
	})
}

// FleetTicker propagates satellite location and velocity on a time interval
func FleetTicker(ticker *time.Ticker, sgp4sats map[string]satellite.Satellite) {

	for t := range ticker.C {
		UpdateSatPos(t, sgp4sats)
	}
}

// GetSatelliteStates pulls the satellite states from the db and converts from json byte to structs
func GetSatelliteStates() map[string]SatelliteState {
	satStates := make(map[string]SatelliteState, 0)
	err := DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB")).Bucket([]byte("FLEET"))
		b.ForEach(func(k, v []byte) error {
			var s SatelliteState
			json.Unmarshal(v, &s)
			satStates[string(k)] = s
			return nil
		})
		return nil
	})
	helpers.PanicErrors(err)
	return satStates
}

// InitSatelliteSGP4 takes satellite state structs and creates a satellite.Satellite object with sgp4 model initialized
func InitSatelliteSGP4(satStates map[string]SatelliteState) map[string]satellite.Satellite {
	sgp4sats := make(map[string]satellite.Satellite, 0)
	for i, sat := range satStates {
		sgp4sats[i] = satellite.TLEToSat(sat.TLELine1, sat.TLELine2, "wgs84")
	}

	return sgp4sats
}
