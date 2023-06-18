package ground

import (
	"fmt"
	"github.com/adrianmo/go-nmea"
	"github.com/igeekinc/go-gps-i2c/pkg/gps"
	"github.com/igeekinc/go-rocket/pkg/core"
	"time"
)

// GPS Rocket Receiver is for rockets that only output GPS data over the wireless serial link
type GPSRocketReceiver struct {
	RocketInfo   *core.RocketInfo
	LastReceived time.Time
	GPS          nmea.GGA // Local GPS info
	port         string
	baudRate     uint
	dataBits     uint
	stopBits     uint
	keepRunning  bool
}

func InitGPSRocketReceiver(rocketInfo *core.RocketInfo, port string, baudRate uint, dataBits uint, stopBits uint) (RocketReceiver, error) {
	rocketReceiver := GPSRocketReceiver{
		RocketInfo:  rocketInfo,
		port:        port,
		baudRate:    baudRate,
		dataBits:    dataBits,
		stopBits:    stopBits,
		keepRunning: false,
	}
	return &rocketReceiver, nil
}

func (recv *GPSRocketReceiver) RocketReceiverLoop() (err error) {
	gpsSerial, err := gps.NewSerialGPSReader(recv.port, recv.baudRate, recv.dataBits, recv.stopBits)
	if err != nil {
		return err
	}

	recv.keepRunning = true

	for recv.keepRunning {
		curFix, err := gpsSerial.NextFix()

		if err == nil {
			readInfo := core.RocketInfo{
				RocketInfoJSON: core.RocketInfoJSON{
					GPS:       curFix,
					Recording: false,
					VideoFile: "",
					Altitude:  curFix.Altitude,
				},
			}
			recv.RocketInfo = &readInfo
			recv.LastReceived = time.Now()

		} else {
			fmt.Printf("Error receiving GPS info from rocket tracker port %s, err %v", recv.port, err)
		}
	}
	return
}

func (recv *GPSRocketReceiver) SendLaunchMode() {
	fmt.Printf("Video capture not supported\n")
}

func (recv *GPSRocketReceiver) UpdateGPS(gpsInfo nmea.GGA) {
	recv.GPS = gpsInfo
}

func (recv *GPSRocketReceiver) GetLocalGPS() nmea.GGA {
	return recv.GPS
}

func (recv *GPSRocketReceiver) GetRocketInfo() core.RocketInfo {
	return *recv.RocketInfo
}

func (recv *GPSRocketReceiver) GetLastReceived() time.Time {
	return recv.LastReceived
}