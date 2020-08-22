package main

import (
	go_rocket "github.com/igeekinc/go-rocket"
	"log"
)

func main() {
	ri := & go_rocket.RocketInfo{}
	gpsReader, err := go_rocket.InitGPSReader(ri, "/dev/ttyS0", 9600, 8, 1)
	if err != nil {
		log.Fatal(err)
	}
	go gpsLoop(gpsReader)
	gpsReporter, err := go_rocket.InitRocketReporter(ri, "/dev/ttyUSB0", 57600, 8, 1)
	reporterLoop(gpsReporter)
}

func gpsLoop(gr go_rocket.GPSReader) {
	err := gr.UpdateFromGPSLoop()
	if err != nil {
		log.Fatal(err)
	}
}

func reporterLoop(gr go_rocket.RocketReporter) {
	err := gr.RocketReporterLoop()
	if err != nil {
		log.Fatal(err)
	}
}