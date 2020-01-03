package go_rocket

import "github.com/adrianmo/go-nmea"

type RocketInfo struct {
	lastGPSTimeStamp nmea.Time
	latitude, longitude float64
}

