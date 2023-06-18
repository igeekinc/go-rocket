package core

import (
	"encoding/json"
	"github.com/adrianmo/go-nmea"
	"periph.io/x/conn/v3/physic"
	"sync"
	"time"
)

type RocketInfoJSON struct {
	Logtime          time.Time
	GPS              nmea.GGA
	Altitude         float64 // Altitude in meters
	XAcc, YAcc, ZAcc physic.Acceleration
	Recording        bool
	VideoFile        string
}
type RocketInfo struct {
	RocketInfoJSON
	lock sync.Mutex
}

func (recv *RocketInfo) UpdateGPS(gps nmea.GGA) {
	recv.lock.Lock()
	defer recv.lock.Unlock()
	recv.GPS = gps
}

func (recv *RocketInfo) UpdateAltitude(temperature float64, pressure float64, altitude float64) {
	recv.Altitude = altitude
}

func (recv *RocketInfo) UpdateAcceleration(x, y, z physic.Acceleration) {
	recv.XAcc = x
	recv.YAcc = y
	recv.ZAcc = z
}

func (recv *RocketInfo) SetRecording(recording bool, videoFile string) {
	recv.lock.Lock()
	defer recv.lock.Unlock()
	recv.Recording = recording
	recv.VideoFile = videoFile
}

func (recv *RocketInfo) MarshalJSON() ([]byte, error) {
	recv.lock.Lock()
	defer recv.lock.Unlock()

	return json.Marshal(recv.RocketInfoJSON)
}

func (recv *RocketInfo) UnmarshalJSON(buf []byte) error {
	if err := json.Unmarshal(buf, &recv.RocketInfoJSON); err != nil {
		return err
	}
	return nil
}
