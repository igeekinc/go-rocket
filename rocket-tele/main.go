package main

import (
	"github.com/igeekinc/go-rocket/pkg/core"
	"github.com/igeekinc/go-rocket/pkg/rocket"
	"log"
)

func main() {
	ri := & core.RocketInfo{}
	gpsReader, err := core.InitGPSReader(ri, "/dev/ttyS0", 9600, 8, 1)
	if err != nil {
		log.Fatal(err)
	}
	go gpsLoop(gpsReader)
	gpsReporter, err := rocket.InitRocketReporter(ri, "/dev/ttyUSB0", 57600, 8, 1)
	reporterLoop(gpsReporter)
}

func gpsLoop(gr core.GPSReader) {
	err := gr.UpdateFromGPSLoop()
	if err != nil {
		log.Fatal(err)
	}
}

func reporterLoop(gr rocket.RocketReporter) {
	err := gr.RocketReporterLoop()
	if err != nil {
		log.Fatal(err)
	}
}