package main

import (
	"github.com/igeekinc/go-rocket/pkg/core"
	"github.com/igeekinc/go-rocket/pkg/ground"
	"log"
	"os"
	"strconv"
)

func main() {
	/*
		app := qtwidgets.NewQApplication(len(os.Args), os.Args, 0)
		btn := qtwidgets.NewQPushButton1("hello qt.go", nil)
		btn.Show()
		app.Exec()
	*/

	tty := os.Args[1]
	baudRate, _ := strconv.Atoi(os.Args[2])
	ri := &core.RocketInfo{}
	rocketReceiver, err := ground.InitRocketReceiver(ri, tty, uint(baudRate), 8, 1)
	if err != nil {
		log.Fatal(err)
	}
	go receiverLoop(rocketReceiver)

	httpServer := ground.NewGroundHTTPServer(".", 8080, rocketReceiver)
	httpServer.Serve()
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
