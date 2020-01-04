package go_rocket

import (
	"encoding/json"
	"github.com/adrianmo/go-nmea"
	"sync"
)

type RocketInfo struct {
	gps nmea.RMC

	lock sync.Mutex
}

func (this *RocketInfo) UpdateGPS(gps nmea.RMC) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.gps = gps
}

func (this *RocketInfo) MarshalJSON() ([] byte, error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	return json.Marshal(this.gps)
}