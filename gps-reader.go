package go_rocket

import (
	"bufio"
	"fmt"
	"github.com/adrianmo/go-nmea"
	"github.com/jacobsa/go-serial/serial"
)

type GPSReader struct {
	rocketInfo  *RocketInfo
	port        string
	baudRate    uint
	dataBits    uint
	stopBits    uint
	keepRunning bool
}

func InitGPSReader(rocketInfo * RocketInfo, port string, baudRate uint, dataBits uint, stopBits uint) (gpsReader GPSReader, err error) {
	gpsReader.rocketInfo = rocketInfo
	gpsReader.port = port
	gpsReader.baudRate = baudRate
	gpsReader.dataBits = dataBits
	gpsReader.stopBits = stopBits
	return gpsReader, nil
}

func (this GPSReader) UpdateFromGPSLoop() (err error) {
	options := serial.OpenOptions{
		PortName:        this.port,
		BaudRate:        this.baudRate,
		DataBits:        this.dataBits,
		StopBits:        this.stopBits,
		MinimumReadSize: 4,
	}

	this.keepRunning = true
	serialPort, err := serial.Open(options)
	if err != nil {
		return
	}
	defer serialPort.Close()
	reader := bufio.NewReader(serialPort)
	scanner := bufio.NewScanner(reader)
	for this.keepRunning {
		for scanner.Scan() {
			var sentence nmea.Sentence
			gpsLine := scanner.Text()
			fmt.Printf("Raw sentence: %v\n", gpsLine)
			sentence, err = nmea.Parse(gpsLine)
			if err == nil {
				if sentence.DataType() == nmea.TypeRMC {
					m := sentence.(nmea.RMC)
					fmt.Printf("Time: %s\n", m.Time)
					fmt.Printf("Validity: %s\n", m.Validity)
					fmt.Printf("Latitude GPS: %s\n", nmea.FormatGPS(m.Latitude))
					fmt.Printf("Latitude DMS: %s\n", nmea.FormatDMS(m.Latitude))
					fmt.Printf("Longitude GPS: %s\n", nmea.FormatGPS(m.Longitude))
					fmt.Printf("Longitude DMS: %s\n", nmea.FormatDMS(m.Longitude))
					fmt.Printf("Speed: %f\n", m.Speed)
					fmt.Printf("Course: %f\n", m.Course)
					fmt.Printf("Date: %s\n", m.Date)
					fmt.Printf("Variation: %f\n", m.Variation)
					this.rocketInfo.UpdateGPS(m)
				}
			}
		}
	}
	return nil
}
