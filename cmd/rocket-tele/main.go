package main

import (
	"github.com/igeekinc/go-rocket/pkg/core"
	"github.com/igeekinc/go-rocket/pkg/rocket"
	"log"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

func main() {
	// Load i2c drivers
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	bus, err := i2creg.Open("/dev/i2c-1")
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	ri := &core.RocketInfo{}

	bmpReader := core.NewBMPReader(bus, ri, 0x77)
	go bmpLoop(bmpReader)

	lsmReader := core.NewLSMReader(bus, ri, 0x19)
	go lsmLoop(lsmReader)

	gpsReader, err := core.InitGPSSerialReader(ri, "/dev/ttyS0", 9600, 8, 1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting gpsLoop")
	go gpsLoop(gpsReader)

	gpsReporter, err := rocket.InitRocketReporter(ri, "/dev/ttyUSB0", 57600, 8, 1)
	reporterLoop(gpsReporter)
}

func gpsLoop(gr *core.GPSReader) {
	err := gr.UpdateFromGPSLoop()
	if err != nil {
		log.Fatal(err)
	}
}

func bmpLoop(bmp *core.BMPReader) {
	err := bmp.UpdateFromBMPLoop()
	if err != nil {
		log.Fatal(err)
	}
}

func lsmLoop(lsm *core.LSMReader) {
	err := lsm.UpdateFromLSMLoop()
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
