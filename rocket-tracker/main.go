package main

import (
	"github.com/igeekinc/go-rocket/pkg/core"
	"github.com/igeekinc/go-rocket/pkg/ground"
	"github.com/igeekinc/go-rocket/pkg/ui"
	"log"
	"os"
	"strconv"
)

func main() {

	tty := os.Args[1]
	baudRate, _ := strconv.Atoi(os.Args[2])

	gpsTTY := os.Args[3]

	ri := &core.RocketInfo{}

	rocketReceiver, err := ground.InitRocketReceiver(ri, tty, uint(baudRate), 8, 1)
	if err != nil {
		log.Fatal(err)
	}
	gpsReader, err := core.InitGPSReader(rocketReceiver, gpsTTY, 9600, 8, 1)
	if err != nil {
		log.Fatal(err)
	}

	go gpsLoop(gpsReader)

	go receiverLoop(rocketReceiver)

	httpServer := ground.NewGroundHTTPServer(".", 8080, rocketReceiver)
	go httpServer.Serve()

	ui.RunRocketTrackerUI(rocketReceiver)
}

func gpsLoop(gr core.GPSReader) {
	err := gr.UpdateFromGPSLoop()
	if err != nil {
		log.Fatal(err)
	}
}

func receiverLoop(rec *ground.RocketReceiver) {
	err := rec.RocketReceiverLoop()
	if err != nil {
		log.Fatal(err)
	}
}
