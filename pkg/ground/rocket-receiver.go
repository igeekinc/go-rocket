package ground

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/adrianmo/go-nmea"
	"github.com/igeekinc/go-rocket/pkg/core"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"time"
)

type RocketReceiver struct {
	RocketInfo   *core.RocketInfo
	LastReceived time.Time
	GPS          nmea.GGA
	port         string
	baudRate     uint
	dataBits     uint
	stopBits     uint
	keepRunning  bool
	serialPort   io.ReadWriteCloser
}

func InitRocketReceiver(rocketInfo *core.RocketInfo, port string, baudRate uint, dataBits uint, stopBits uint) (*RocketReceiver, error) {
	rocketReceiver := RocketReceiver{
		RocketInfo:  rocketInfo,
		port:        port,
		baudRate:    baudRate,
		dataBits:    dataBits,
		stopBits:    stopBits,
		keepRunning: false,
		serialPort:  nil,
	}
	return &rocketReceiver, nil
}

func (recv *RocketReceiver) RocketReceiverLoop() (err error) {
	options := serial.OpenOptions{
		PortName:        recv.port,
		BaudRate:        recv.baudRate,
		DataBits:        recv.dataBits,
		StopBits:        recv.stopBits,
		MinimumReadSize: 4,
	}

	serialPort, err := serial.Open(options)
	if err != nil {
		return
	}
	recv.serialPort = serialPort
	defer recv.serialPort.Close()

	serialPortReader := bufio.NewReaderSize(recv.serialPort, 16*1024)
	recv.keepRunning = true

	for recv.keepRunning {
		jsonBytes, err := serialPortReader.ReadBytes('\n')
		if err == nil {
			jsonStr := string(jsonBytes)
			var readInfo core.RocketInfo
			err = json.Unmarshal(jsonBytes, &readInfo)
			if err == nil {
				fmt.Printf("Updating from %s\n", jsonStr)
				*recv.RocketInfo = readInfo
				recv.LastReceived = time.Now()
				fmt.Printf("RocketInfo.logtime=%v, now = %v\n", recv.RocketInfo.Logtime, time.Now())
			} else {
				fmt.Printf("Could not unmarshal data from serial port %s, data = '%s', err = %v\n",
					recv.port, jsonStr, err)
			}
		} else {
			fmt.Printf("Got error reading from serial port %s, err %v\n", recv.port, err)
		}
	}
	return
}

func (recv *RocketReceiver) SendLaunchMode() {
	fmt.Printf("Sending video cmd to rocket\n")
	written, err := recv.serialPort.Write([]byte("\n\nV\n"))
	if err != nil {
		fmt.Printf("Wrote %d bytes, got error %v\n", written, err)
	} else {
		fmt.Printf("Wrote %d bytes\n", written)
	}
}

func (recv *RocketReceiver) UpdateGPS(gpsInfo nmea.GGA) {
	recv.GPS = gpsInfo
}
