package core

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

type RocketInfoLog interface {
	BaseAGL() float64
	NumRecords() int
	Record(recordNum int) RocketInfo
	BaroHeight(recordNum int) float64
	BaroHeightAGL(recordNum int) float64
	MaxBaroHeight() (float64, RocketInfo)
	MaxBaroHeightAGL() (float64, RocketInfo)
}

type rocketInfoLog struct {
	logFilePath   string
	records       []RocketInfo
	baseAGL       float64
	baroGPSOffset float64
	aglOffset     float64
}

func ReadRocketInfoLog(logFilePath string) (RocketInfoLog, error) {
	ril := rocketInfoLog{
		logFilePath: logFilePath,
	}
	logFile, err := os.Open(logFilePath)
	if err != nil {
		log.Fatalf("Could not open log file %s:%v", logFilePath, err)
	}

	logFileBufIO := bufio.NewReader(logFile)

	records := []RocketInfo{}
	lineNum := 0
	for true {
		curLine, err := logFileBufIO.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				log.Printf("Error reading line %d, err: %v\n", lineNum, err)
			}
		}
		curRocketInfo := RocketInfo{}
		err = json.Unmarshal([]byte(curLine), &curRocketInfo)
		if err != nil {
			log.Printf("Could not parse line %d:'%s' err:%v\n", lineNum, curLine, err)
		}
		lineNum++
		records = append(records, curRocketInfo)
	}

	if len(records) == 0 {
		log.Printf("Did not read any lines from log file %s\n", logFilePath)
		return nil, err
	}

	ril.records = records

	ril.baseAGL = records[0].GPS.Altitude

	ril.baroGPSOffset = records[0].GPS.Altitude - records[0].Altitude
	ril.aglOffset = (0 - records[0].Altitude) - ril.baroGPSOffset
	return &ril, nil
}

func (recv *rocketInfoLog) BaseAGL() float64 {
	return recv.baseAGL
}

func (recv *rocketInfoLog) NumRecords() int {
	return len(recv.records)
}

func (recv *rocketInfoLog) BaroHeight(recordNum int) float64 {
	return recv.records[recordNum].Altitude + recv.baroGPSOffset
}

func (recv *rocketInfoLog) BaroHeightAGL(recordNum int) float64 {
	return recv.records[recordNum].Altitude + recv.baroGPSOffset + recv.aglOffset
}

func (recv *rocketInfoLog) MaxBaroHeightAGL() (float64, RocketInfo) {
	maxBaroHeight, info := recv.MaxBaroHeight()
	maxBaroHeightAGL := maxBaroHeight - recv.aglOffset
	return maxBaroHeightAGL, info
}

func (recv *rocketInfoLog) MaxBaroHeight() (float64, RocketInfo) {
	var maxBaroHeight float64
	maxHeightRecordNum := 0
	for recordNum := 0; recordNum < recv.NumRecords(); recordNum++ {
		if recv.BaroHeight(recordNum) > maxBaroHeight {
			maxBaroHeight = recv.BaroHeight(recordNum)
			maxHeightRecordNum = recordNum
		}
	}
	return maxBaroHeight, recv.records[maxHeightRecordNum]
}

func (recv *rocketInfoLog) Record(recordNum int) RocketInfo {
	return recv.records[recordNum]
}
