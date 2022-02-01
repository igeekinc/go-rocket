package core

import (
	"encoding/json"
	"github.com/adrianmo/go-nmea"
	"sync"
)

type RocketInfo struct {
	GPS nmea.RMC

	lock sync.Mutex
}

func (this *RocketInfo) UpdateGPS(gps nmea.RMC) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.GPS = gps
}

func (this *RocketInfo) MarshalJSON() ([] byte, error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	return json.Marshal(this.GPS)
}


func (this *RocketInfo) UnmarshalJSON(buf []byte) error {
	if err := json.Unmarshal(buf, &this.GPS); err != nil {
		return err
	}
	return nil
}