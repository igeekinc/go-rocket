package core

import (
	"log"
	"math"
	i2c2 "periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
	"time"
)

type AltitudeTracker interface {
	UpdateAltitude(temperature float64, pressure float64, altitude float64)
}

type BMPReader struct {
	altitudeTracker AltitudeTracker
	i2cBus          i2c2.BusCloser
	bmpI2CAddr      uint16
}

func NewBMPReader(i2cBus i2c2.BusCloser, altitudeTracker AltitudeTracker, bmpI2CAddr uint16) *BMPReader {
	return &BMPReader{
		altitudeTracker: altitudeTracker,
		i2cBus:          i2cBus,
		bmpI2CAddr:      bmpI2CAddr,
	}
}

func (recv *BMPReader) UpdateFromBMPLoop() (err error) {
	dev, err := bmxx80.NewI2C(recv.i2cBus, recv.bmpI2CAddr, &bmxx80.DefaultOpts)
	if err != nil {
		return err
	}
	for {
		var env physic.Env
		if err = dev.Sense(&env); err != nil {
			log.Fatal(err)
		}
		seaLevelPA := 101325.0 // Sea level in Pascals
		pressurePA := float64(env.Pressure / physic.Pascal)
		altitude := 44330 * (1 - math.Pow(pressurePA/seaLevelPA, 1/5.255))
		// Round up to 2 decimals after point
		altitudeRounded := float64(int(altitude*100)) / 100
		recv.altitudeTracker.UpdateAltitude(env.Temperature.Celsius(), float64(env.Pressure)/float64(physic.Pascal), altitudeRounded)
		/*fmt.Printf("%8s %10s %9s altitude %.2f pressurePA %.2f seaLevelPA %.2f\n", env.Temperature, env.Pressure, env.Humidity, altitudeRounded,
		pressurePA, seaLevelPA)*/
		time.Sleep(100 * time.Millisecond) // Update 10 times per second
		//time.Sleep(time.Second)
	}
}
