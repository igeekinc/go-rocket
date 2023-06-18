package core

import (
	"github.com/bskari/go-lsm303/pkg/lsm303"
	"github.com/pkg/errors"
	i2c2 "periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/physic"
	"time"
)

type AccelerationTracker interface {
	UpdateAcceleration(x, y, z physic.Acceleration)
}

type LSMReader struct {
	accelerationTracker AccelerationTracker
	i2cBus              i2c2.BusCloser
	bmpI2CAddr          uint16
}

func NewLSMReader(i2cBus i2c2.BusCloser, accelerationTracker AccelerationTracker, bmpI2CAddr uint16) *LSMReader {
	return &LSMReader{
		accelerationTracker: accelerationTracker,
		i2cBus:              i2cBus,
		bmpI2CAddr:          bmpI2CAddr,
	}
}

func (recv *LSMReader) UpdateFromLSMLoop() (err error) {
	accOpts := lsm303.AccelerometerOpts{
		Range: lsm303.ACCELEROMETER_RANGE_16G,
		Mode:  lsm303.ACCELEROMETER_MODE_NORMAL,
	}
	dev, err := lsm303.NewAccelerometer(recv.i2cBus, &accOpts)

	if err != nil {
		return errors.Wrap(err, "Could not connect to accelerometer")
	}

	for {
		x, y, z, err := dev.Sense()
		if err == nil {
			recv.accelerationTracker.UpdateAcceleration(x, y, z)
		}
		time.Sleep(100 * time.Millisecond) // Update 10 times per second
	}
}
