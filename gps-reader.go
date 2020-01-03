package go_rocket

import (
	"bufio"
	"fmt"
	"github.com/adrianmo/go-nmea"
	"github.com/jacobsa/go-serial/serial"
)
type GPSReader struct {
	port string
	baudRate uint
	dataBits uint
	stopBits uint
	keepRunning bool
}

func InitGPSReader(port string, baudRate uint, dataBits uint, stopBits uint) (gpsReader GPSReader, err error) {
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
			if sentence.DataType() == nmea.TypeRMC {

			}
		}
	}
	return nil
}