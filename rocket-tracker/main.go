package main

import (
	"fmt"
	go_rocket "github.com/igeekinc/go-rocket"
	"log"
	"math"
	"time"
)

func main() {
	/*
		app := qtwidgets.NewQApplication(len(os.Args), os.Args, 0)
		btn := qtwidgets.NewQPushButton1("hello qt.go", nil)
		btn.Show()
		app.Exec()
	*/

	ri := &go_rocket.RocketInfo{}
	rocketReceiver, err := go_rocket.InitRocketReceiver(ri, "/dev/tty.SLAB_USBtoUART", 57600, 8, 1)
	if err != nil {
		log.Fatal(err)
	}
	go receiverLoop(rocketReceiver)

	ourPos := &go_rocket.RocketInfo{}
	gpsReader, err := go_rocket.InitGPSReader(ourPos, "/dev/tty.usbmodem14222101", 9600, 8, 1)
	if err != nil {
		log.Fatal(err)
	}

	go gpsLoop(gpsReader)

	for true {
		distance := Distance(ourPos.GPS.Latitude, ourPos.GPS.Longitude, ri.GPS.Latitude, ri.GPS.Longitude)
		fmt.Printf("Distance to rocket is %f meters\n", distance)
		time.Sleep(1 * time.Second)

	}
}

func gpsLoop(gr go_rocket.GPSReader) {
	err := gr.UpdateFromGPSLoop()
	if err != nil {
		log.Fatal(err)
	}
}

func receiverLoop(rec go_rocket.RocketReceiver) {
	err := rec.RocketReceiverLoop()
	if err != nil {
		log.Fatal(err)
	}
}

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}
