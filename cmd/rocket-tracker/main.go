package main

import (
	"flag"
	"fmt"
	"github.com/igeekinc/go-rocket/pkg/core"
	"github.com/igeekinc/go-rocket/pkg/ground"
	"github.com/igeekinc/go-rocket/pkg/ui"
	"log"
	"os"
)

func main() {
	goRocketTTY := flag.String("go-rocket-tty", "", "TTY with input data from a go-rocket rocket-tele stream")
	gpsTTY := flag.String("gps-tty", "", "TTY connected to tracker console/handheld's GPS")
	gpsTrackerTTY := flag.String("gps-tracker-tty", "", "TTY connected to rocket GPS only tracker")
	baudRate := flag.Int("baud-rate", 57600, "Baud rate from rocket")

	flag.Parse()

	if gpsTTY == nil || *gpsTTY == "" {
		fmt.Println("Must provide a GPS TTY")
		flag.Usage()
		os.Exit(1)
	}

	if (goRocketTTY == nil || *goRocketTTY == "") && (gpsTrackerTTY == nil || *gpsTrackerTTY == "") {
		fmt.Println("Must provide a Go rocket tracker TTY or a GPS tracker TTY")
		flag.Usage()
		os.Exit(1)
	}

	if goRocketTTY != nil && *goRocketTTY != "" && gpsTrackerTTY != nil && *gpsTrackerTTY != "" {
		fmt.Println("Cannot provide both a Go rocket tracker TTY and a GPS tracker TTY")
		flag.Usage()
		os.Exit(1)
	}
	ri := &core.RocketInfo{}

	var err error
	var rocketReceiver ground.RocketReceiver
	if goRocketTTY != nil && *goRocketTTY != "" {
		fmt.Printf("Initialize remote go-rocket tracker on %s\n", *goRocketTTY)
		rocketReceiver, err = ground.InitGoRocketReceiver(ri, *goRocketTTY, uint(*baudRate), 8, 1)
		if err != nil {
			log.Fatal(err)
		}
	}

	if gpsTrackerTTY != nil && *gpsTrackerTTY != "" {
		fmt.Printf("Initialize GPS Only tracker on %s\n", *gpsTrackerTTY)

		rocketReceiver, err = ground.InitGPSRocketReceiver(ri, *gpsTrackerTTY, uint(*baudRate), 8, 1)
		if err != nil {
			log.Fatal(err)
		}
	}

	gpsReader, err := core.InitGPSSerialReader(rocketReceiver, *gpsTTY, 9600, 8, 1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Calling gpsLoop")
	go gpsLoop(gpsReader)

	go receiverLoop(rocketReceiver)

	httpServer := ground.NewGroundHTTPServer(".", 8080, rocketReceiver)
	go httpServer.Serve()

	ui.RunRocketTrackerUI(rocketReceiver)
}

func gpsLoop(gr *core.GPSReader) {
	log.Println("Calling UpdateFromGPSLoop")
	err := gr.UpdateFromGPSLoop()
	if err != nil {
		log.Fatal(err)
	}
}

func receiverLoop(rec ground.RocketReceiver) {
	err := rec.RocketReceiverLoop()
	if err != nil {
		log.Fatal(err)
	}
}
