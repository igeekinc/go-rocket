package go_rocket

import (
	"encoding/json"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"time"
)

type GPSReporter struct {
	rocketInfo *RocketInfo
	port        string
	baudRate    uint
	dataBits    uint
	stopBits    uint
	keepRunning bool
}

func InitGPSReporter(rocketInfo *RocketInfo, port string, baudRate uint, dataBits uint, stopBits uint) (gpsReporter GPSReporter, err error) {
	gpsReporter.rocketInfo = rocketInfo
	gpsReporter.port = port
	gpsReporter.baudRate = baudRate
	gpsReporter.dataBits = dataBits
	gpsReporter.stopBits = stopBits
	return gpsReporter, nil
}

func (this *GPSReporter) GPSReporterLoop() (err error) {
	options := serial.OpenOptions{
		PortName:        this.port,
		BaudRate:        this.baudRate,
		DataBits:        this.dataBits,
		StopBits:        this.stopBits,
		MinimumReadSize: 4,
	}


	serialPort, err := serial.Open(options)
	if err != nil {
		return
	}
	defer serialPort.Close()

	this.keepRunning = true

	for this.keepRunning {
		rj, err := json.Marshal(this.rocketInfo)
		if err != nil {
			return err
		}
		str := string(rj)
		fmt.Println(str)
		bytesWritten, err := serialPort.Write(rj)
		if err != nil {
			fmt.Printf("Wrote %d bytes and got error %v", bytesWritten, err)
		}
		time.Sleep(1 * time.Second)

	}
	return
}

