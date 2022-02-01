package go_rocket

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/igeekinc/go-rocket/pkg/go-rocket-core"
	"github.com/jacobsa/go-serial/serial"
)

type RocketReceiver struct {
	rocketInfo  *go_rocket_core.RocketInfo
	port        string
	baudRate    uint
	dataBits    uint
	stopBits    uint
	keepRunning bool
}

func InitRocketReceiver(rocketInfo *go_rocket_core.RocketInfo, port string, baudRate uint, dataBits uint, stopBits uint) (rocketReceiver RocketReceiver, err error) {
	rocketReceiver.rocketInfo = rocketInfo
	rocketReceiver.port = port
	rocketReceiver.baudRate = baudRate
	rocketReceiver.dataBits = dataBits
	rocketReceiver.stopBits = stopBits
	return rocketReceiver, nil
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
	defer serialPort.Close()

	serialPortReader := bufio.NewReaderSize(serialPort, 16*1024)
	this.keepRunning = true

	for this.keepRunning {
		jsonBytes, err := serialPortReader.ReadBytes('\n')
		if (err == nil) {
			jsonStr := string(jsonBytes)
			var readInfo go_rocket_core.RocketInfo
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
