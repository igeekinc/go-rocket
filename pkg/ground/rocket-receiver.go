package ground

import (
	"github.com/adrianmo/go-nmea"
	"github.com/igeekinc/go-rocket/pkg/core"
	"time"
)

type RocketReceiver interface {
	RocketReceiverLoop() (err error)
	SendLaunchMode()
	UpdateGPS(gpsInfo nmea.GGA)
	GetLocalGPS() nmea.GGA
	GetRocketInfo() core.RocketInfo
	GetLastReceived() time.Time
}
