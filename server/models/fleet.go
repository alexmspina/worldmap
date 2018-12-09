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

	"github.com/joshuaferrara/go-satellite"

	"github.com/boltdb/bolt"
)

// Fleet struct of all satellite states derived from the tle and beamplan
type Fleet struct {
	Satellites map[string]SatelliteState `json:"satellites"`
}

// SatelliteState struct that models the static info of a satellite including tle lines and current beamplan
type SatelliteState struct {
	TLELine1 string                     `json:"tleLine1"`
	TLELine2 string                     `json:"tleLine2"`
	Missions map[string]BeamplanMission `json:"missions"`
}

// BeamplanMission struct modeling beamplan mission
type BeamplanMission struct {
	MissionConfig          string                    `json:"missionConfig"`
	GatewayTargetID        string                    `json:"gatewayTargetID"`
	GatewayOBAntID         string                    `json:"gatewayOBAntID"`
	GatewayPointingMaxTime string                    `json:"gatewayPointingMaxTime"`
	Targets                map[string]BeamProperties `json:"targets"`
}

// BeamProperties struct modeling individual beam settings
type BeamProperties struct {
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

// GetTLES creats a map of tles with map of tle lines
func GetTLES(tle string) map[string]map[string]string {
	// get pointer to file of tle
	tlefile, err := os.Open(tle)
	PanicErrors(err)

	// read file into byte slices
	tlereader := io.Reader(tlefile)
	tlelines, err := ioutil.ReadAll(tlereader)
	PanicErrors(err)

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
func GetBeamplan(tlemap map[string]map[string]string, bpfiles map[string]string, db *bolt.DB) {
	spares := []string{"M002", "M004", "M005"}
	activenotb3 := []string{"M001", "M003", "M006", "M007", "M008", "M009", "M010", "M011", "M012"}
	activeb3 := []string{"M013", "M014", "M015", "M016"}

	for sat, tle := range tlemap {
		switch true {
		case StringInSlice(sat, spares):
			bpfile := bpfiles["SPARE"]
			satstate := BuildSatelliteState(bpfile, tle, sat)
			FillFleetBucket(sat, satstate, db)
		case StringInSlice(sat, activenotb3):
			switch sat {
			case "M001":
				bpfile := bpfiles["M001"]
				satstate := BuildSatelliteState(bpfile, tle, sat)
				FillFleetBucket(sat, satstate, db)
			default:
				bpfile := bpfiles["ACTIVE"]
				satstate := BuildSatelliteState(bpfile, tle, sat)
				FillFleetBucket(sat, satstate, db)
			}
		case StringInSlice(sat, activeb3):
			switch sat {
			case "M013":
				bpfile := bpfiles["M013"]
				satstate := BuildSatelliteState(bpfile, tle, sat)
				FillFleetBucket(sat, satstate, db)
			default:
				bpfile := bpfiles["B3"]
				satstate := BuildSatelliteState(bpfile, tle, sat)
				FillFleetBucket(sat, satstate, db)
			}
		}
	}
}

// FillFleetBucket initializes satellites from tle
func FillFleetBucket(satid string, satstate SatelliteState, db *bolt.DB) error {
	satstatebytes, err := json.MarshalIndent(satstate, "", "\t")
	if err != nil {
		panic(err)
	}
	satstatebytes = bytes.Replace(satstatebytes, []byte("\\u0026"), []byte("&"), -1)
	satstatebytes = bytes.Trim(satstatebytes, "\r")

	err = db.Update(func(tx *bolt.Tx) error {
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
			PanicErrors(err)
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

	msnsmap := make(map[string]BeamplanMission, 0)
	for i, beams := range msns {
		tgts := make(map[string]BeamProperties, 0)
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

			tgts[tgtid] = bmprops
		}

		bpmsn := BeamplanMission{
			MissionConfig:          msnconfig,
			GatewayTargetID:        gwtgtid,
			GatewayOBAntID:         gwobantid,
			GatewayPointingMaxTime: gwpntmxtime,
			Targets:                tgts,
		}

		msnsmap[i] = bpmsn
	}

	satstate := SatelliteState{
		TLELine1: tle["firstline"],
		TLELine2: tle["secondline"],
		Missions: msnsmap,
	}

	return satstate
}

// FleetTicker propagates satellite location and velocity on a time interval
func FleetTicker(db *bolt.DB) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		// select {
		// case t := <-ticker.C:
		// 	fmt.Println("Current time: ", t)
		// }
		t := <-ticker.C
		var v []byte

		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("DB")).Bucket([]byte("FLEET"))
			// b.ForEach(func(k, v []byte) error {
			// 	fmt.Println(string(k), string(v))
			// 	return nil
			// })
			// return nil
			v = b.Get([]byte("M003"))
			return nil
		})
		PanicErrors(err)

		var s SatelliteState
		json.Unmarshal(v, &s)
		// parsedtle := satellite.ParseTLE(s.TLELine1, s.TLELine2, "wgs84")
		utc := t.UTC()
		sat := satellite.TLEToSat(s.TLELine1, s.TLELine2, "wgs84")
		y, m, d := utc.Date()
		h, min, sec := utc.Clock()
		gmst := satellite.GSTimeFromDate(y, int(m), d, h, min, sec)
		pos, _ := satellite.Propagate(sat, y, int(m), d, h, min, sec)
		_, _, latlng := satellite.ECIToLLA(pos, gmst)

		// fmt.Println(parsedtle)
		fmt.Printf("Latlng: %v", latlng)
		fmt.Println("Current time: ", t)
	}
}
