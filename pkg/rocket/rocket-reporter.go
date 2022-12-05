package rocket

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/igeekinc/go-rocket/pkg/core"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type RocketReporter struct {
	rocketInfo  *core.RocketInfo
	port        string
	baudRate    uint
	dataBits    uint
	stopBits    uint
	keepRunning bool
}

func InitRocketReporter(rocketInfo *core.RocketInfo, port string, baudRate uint, dataBits uint, stopBits uint) (rocketReporter RocketReporter, err error) {
	rocketReporter.rocketInfo = rocketInfo
	rocketReporter.port = port
	rocketReporter.baudRate = baudRate
	rocketReporter.dataBits = dataBits
	rocketReporter.stopBits = stopBits
	return rocketReporter, nil
}

func (recv *RocketReporter) GetInfo() core.RocketInfo {
	retInfo := *recv.rocketInfo
	retInfo.Logtime = time.Now()
	return retInfo
}

func (recv *RocketReporter) RocketReporterLoop() (err error) {
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
	defer serialPort.Close()

	go recv.videoStarter(serialPort)
	recv.keepRunning = true

	go recv.videoStarter(serialPort)

	for recv.keepRunning {
		rj, err := json.Marshal(recv.GetInfo())
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

const flightTime = 10 * time.Minute

//const flightTime = 10 * time.Second

func (recv *RocketReporter) videoStarter(serialPort io.Reader) error {
	fmt.Println("====================")
	fmt.Println("videoStarter running")
	fmt.Println("====================")
	scanner := bufio.NewScanner(serialPort)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println("====================")
		fmt.Println(text)
		fmt.Println("====================")
		if text == "V" {
			videoFile, logFile, err := recv.nextFlightFiles()
			if err != nil {
				return err
			}
			fmt.Printf("Starting logfile = %s\n", logFile)
			logFinished := make(chan bool)
			go recv.logFlight(logFinished, logFile, flightTime)

			recv.rocketInfo.SetRecording(true, videoFile)
			video := exec.Command("/usr/bin/libcamera-vid", "-t", strconv.FormatInt(flightTime.Milliseconds(), 10),
				"--nopreview", "-o", videoFile)
			fmt.Printf("Starting video, cmd = %v\n", video)
			video.Run()
			fmt.Printf("Waiting for logger to finish\n")
			<-logFinished
			recv.rocketInfo.SetRecording(false, videoFile)
			fmt.Printf("Finished flight mode\n")
		}
	}
	return nil
}

const logSleepTime = time.Millisecond * 10

func (recv *RocketReporter) logFlight(finished chan bool, logFile string, duration time.Duration) {
	endTime := time.Now().Add(duration)
	logFileOut, err := os.Create(logFile)
	if err == nil {
		for time.Now().Before(endTime) {
			rj, err := json.Marshal(recv.GetInfo())
			if err == nil {
				rj = append(rj, '\n')
				logFileOut.Write(rj)
			}
			time.Sleep(logSleepTime)
		}
	}
	fmt.Printf("logFlight finished\n")
	logFileOut.Close()
	finished <- true
}

const kFlightsDir = "/flights"

func (recv *RocketReporter) nextFlightFiles() (vidFile, logFile string, err error) {
	info, err := ioutil.ReadDir(kFlightsDir)
	if err != nil {
		panic(err) // If we're having errors reading the video dir, just bomb out
	}
	nextFlightNum := 0
	for _, curFile := range info {
		var curFlightNum int
		n, err := fmt.Sscanf(curFile.Name(), "flightdir%d", &curFlightNum)
		if err != nil {
			panic(err)
		}
		if n == 1 {
			if curFlightNum > nextFlightNum {
				nextFlightNum = curFlightNum
			}
		}
	}
	nextFlightNum++
	flightDir := fmt.Sprintf("%s/flightdir%d", kFlightsDir, nextFlightNum)
	err = os.MkdirAll(flightDir, 0777)
	if err == nil {
		vidFile = fmt.Sprintf("%s/vid%d.mov", flightDir, nextFlightNum)
		logFile = fmt.Sprintf("%s/log%d.json", flightDir, nextFlightNum)
	}
	return
}
