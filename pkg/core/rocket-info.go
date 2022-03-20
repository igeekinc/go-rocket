package core

import (
	"encoding/json"
	"github.com/adrianmo/go-nmea"
	"sync"
)

type RocketInfoJSON struct {
	GPS       nmea.GGA
	Altitude  float64 // Altitude in meters
	Recording bool
	VideoFile string
}
type RocketInfo struct {
	RocketInfoJSON
	lock sync.Mutex
}

func (this *RocketInfo) UpdateGPS(gps nmea.GGA) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.GPS = gps
}

func (this *RocketInfo) UpdateAltitude(temperature float64, pressure float64, altitude float64) {
	this.Altitude = altitude
}

func (this *RocketInfo) SetRecording(recording bool, videoFile string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.Recording = recording
	this.VideoFile = videoFile
}

func (this *RocketInfo) MarshalJSON() ([]byte, error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	return json.Marshal(this.RocketInfoJSON)
}

func (this *RocketInfo) UnmarshalJSON(buf []byte) error {
	if err := json.Unmarshal(buf, &this.RocketInfoJSON); err != nil {
		return err
	}
	return nil
}
