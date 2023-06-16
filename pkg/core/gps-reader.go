package core

import (
	"fmt"
	"github.com/adrianmo/go-nmea"
	"github.com/igeekinc/go-gps-i2c/pkg/gps"
	"log"
	"periph.io/x/conn/v3/i2c"
	"time"
)

type GPSReader struct {
	// A consumer of the GPS info that will be periodically updated - TODO - convert over to using a channel or something
	// more idiomatic
	gpsTracker GPSTracker
	// The GPS device itself
	gps         gps.GPS
	port        string
	baudRate    uint
	dataBits    uint
	stopBits    uint
	keepRunning bool
}

type GPSTracker interface {
	UpdateGPS(gpsInfo nmea.GGA)
}

func InitGPSSerialReader(gpsTracker GPSTracker, port string, baudRate uint, dataBits uint, stopBits uint) (*GPSReader, error) {
	gpsSerial, err := gps.NewSerialGPSReader(port, baudRate, dataBits, stopBits)
	if err != nil {
		return nil, err
	}
	gpsReader := GPSReader{
		gpsTracker:  gpsTracker,
		gps:         gpsSerial,
		keepRunning: true,
	}
	return &gpsReader, nil
}

func InitGPSI2CReader(gpsTracker GPSTracker, bus i2c.BusCloser, opts *gps.Opts) (*GPSReader, error) {
	gpsI2C, err := gps.NewI2CGPS(bus, opts)
	if err != nil {
		return nil, err
	}
	gpsReader := GPSReader{
		gpsTracker:  gpsTracker,
		gps:         gpsI2C,
		keepRunning: true,
	}
	return &gpsReader, nil
}

func (recv *GPSReader) UpdateFromGPSLoop() (err error) {
	for recv.keepRunning {
		log.Printf("Getting GPS local fix")
		g, err := recv.gps.NextFix()
		if err != nil {
			log.Printf("GPS error - err = %v\n", err)
			time.Sleep(time.Second)
			continue
		}
		fmt.Printf("Time: %s\n", g.Time)
		fmt.Printf("Fix quality: %s\n", g.FixQuality)
		fmt.Printf("Latitude GPS: %s\n", nmea.FormatGPS(g.Latitude))
		fmt.Printf("Longitude GPS: %s\n", nmea.FormatGPS(g.Longitude))
		lastGSVTime, lastGSV := recv.gps.LastGSV()
		fmt.Printf("Last GSV time:%v\n", lastGSVTime)
		fmt.Printf("Num satellites: %d\n", lastGSV.NumberSVsInView)
		fmt.Printf("Satellite info: %v\n", lastGSV.Info)
		recv.gpsTracker.UpdateGPS(g)
	}
	return nil
}
