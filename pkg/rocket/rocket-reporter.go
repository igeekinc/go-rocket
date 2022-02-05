package rocket

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/igeekinc/go-rocket/pkg/core"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"os/exec"
	"time"
)

type RocketReporter struct {
	rocketInfo *core.RocketInfo
	port        string
	baudRate    uint
	dataBits    uint
	stopBits    uint
	keepRunning bool
	video       bool
}

func InitRocketReporter(rocketInfo *core.RocketInfo, port string, baudRate uint, dataBits uint, stopBits uint) (rocketReporter RocketReporter, err error) {
	rocketReporter.rocketInfo = rocketInfo
	rocketReporter.port = port
	rocketReporter.baudRate = baudRate
	rocketReporter.dataBits = dataBits
	rocketReporter.stopBits = stopBits
	return rocketReporter, nil
}

func (this *RocketReporter) RocketReporterLoop() (err error) {
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
		rj = append(rj, '\r')
		rj = append(rj, '\n')
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

func (this *RocketReporter) videoStarter(serialPort io.Reader) {
	scanner := bufio.NewScanner(serialPort)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println(text)
		if text == "V" {
			this.video = true
			video := exec.Command("/usr/bin/raspivid", "--timeout", "6000000",
				"-o", this.nextVidFile())
			video.Run()
			this.video = false
		}
	}

}

func (this *RocketReporter) nextVidFile() string {
	return "vid.mov"
}
