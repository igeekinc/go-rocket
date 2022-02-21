package ground

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/igeekinc/go-rocket/pkg/core"
	"github.com/jacobsa/go-serial/serial"
	"io"
)

type RocketReceiver struct {
	rocketInfo  *core.RocketInfo
	port        string
	baudRate    uint
	dataBits    uint
	stopBits    uint
	keepRunning bool
	serialPort  io.ReadWriteCloser
}

func InitRocketReceiver(rocketInfo *core.RocketInfo, port string, baudRate uint, dataBits uint, stopBits uint) (*RocketReceiver, error) {
	rocketReceiver := RocketReceiver{
		rocketInfo:  rocketInfo,
		port:        port,
		baudRate:    baudRate,
		dataBits:    dataBits,
		stopBits:    stopBits,
		keepRunning: false,
		serialPort:  nil,
	}
	return &rocketReceiver, nil
}

func (this *RocketReceiver) RocketReceiverLoop() (err error) {
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
	this.serialPort = serialPort
	defer this.serialPort.Close()

	serialPortReader := bufio.NewReaderSize(this.serialPort, 16*1024)
	this.keepRunning = true

	for this.keepRunning {
		jsonBytes, err := serialPortReader.ReadBytes('\n')
		if err == nil {
			jsonStr := string(jsonBytes)
			var readInfo core.RocketInfo
			err = json.Unmarshal(jsonBytes, &readInfo)
			if err == nil {
				fmt.Printf("Updating from %s\n", jsonStr)
				*this.rocketInfo = readInfo
			} else {
				fmt.Printf("Could not unmarshal data from serial port %s, data = '%s', err = %v\n",
					this.port, jsonStr, err)
			}
		} else {
			fmt.Printf("Got error reading from serial port %s, err %v\n", this.port, err)
		}
	}
	return
}
