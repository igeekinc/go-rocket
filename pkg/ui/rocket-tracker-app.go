package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/adrianmo/go-nmea"
	"github.com/igeekinc/go-rocket/pkg/core"
	"github.com/igeekinc/go-rocket/pkg/ground"
	"time"
)

func RunRocketTrackerUI(rocketReceiver *ground.RocketReceiver) {
	a := app.New()
	w := a.NewWindow("Rocket Tracker")
	//w.SetFullScreen(true)
	//w.SetContent(widget.NewLabel("Hello World!"))
	useGrid := false
	kRocketPosLabel := widget.NewLabel("Rckt:")
	rocketPosLabel := widget.NewLabel("")
	kOurPosLabel := widget.NewLabel("Us:")
	ourPosLabel := widget.NewLabel("")
	kDistanceLabel := widget.NewLabel("Dist:")
	distanceLabel := widget.NewLabel("")
	kAccelerationLabel := widget.NewLabel("Acc:")
	accelerationLabel := widget.NewLabel("")
	kLaunchModeLabel := widget.NewLabel("Mode: ")
	launchModeLabel := widget.NewLabel("ready")

	rtu := rocketTrackerUI{
		rocketReceiver:       rocketReceiver,
		gpsPosLabel:          rocketPosLabel,
		ourPosLabel:          ourPosLabel,
		distanceLabel:        distanceLabel,
		accelerationLabel:    accelerationLabel,
		launchModeStateLabel: launchModeLabel,
	}
	launchButton := widget.NewButton("Launch", func() {
		fmt.Println("Launch pressed")
		rocketReceiver.SendLaunchMode()
	})

	var infoCanvas fyne.CanvasObject
	if useGrid {
		infoCanvas = container.New(layout.NewGridLayout(2),
			kRocketPosLabel,
			rocketPosLabel,
			kOurPosLabel,
			ourPosLabel,
			kDistanceLabel,
			distanceLabel,
			kAccelerationLabel,
			accelerationLabel,
			kLaunchModeLabel,
			launchModeLabel,
			launchButton)

	} else {
		labelContainer := container.New(layout.NewVBoxLayout(),
			kRocketPosLabel,
			kOurPosLabel,
			kDistanceLabel,
			kAccelerationLabel,
			kLaunchModeLabel)
		infoContainer := container.New(layout.NewVBoxLayout(),
			rocketPosLabel,
			ourPosLabel,
			distanceLabel,
			accelerationLabel,
			launchModeLabel)
		comboInfoContainer := container.New(layout.NewHBoxLayout(),
			labelContainer, infoContainer)
		infoCanvas = container.New(layout.NewVBoxLayout(),
			comboInfoContainer, launchButton)
	}
	w.SetContent(infoCanvas)

	go rtu.updateLoop()
	w.ShowAndRun()
}

type rocketTrackerUI struct {
	rocketReceiver       *ground.RocketReceiver
	gpsPosLabel          *widget.Label
	ourPosLabel          *widget.Label
	distanceLabel        *widget.Label
	accelerationLabel    *widget.Label
	launchModeStateLabel *widget.Label
}

func (recv *rocketTrackerUI) updateLoop() {
	for {
		updateLabelWithGPS(recv.gpsPosLabel, recv.rocketReceiver.RocketInfo.GPS, recv.rocketReceiver.RocketInfo.Altitude)
		updateLabelWithGPS(recv.ourPosLabel, recv.rocketReceiver.GPS, recv.rocketReceiver.GPS.Altitude)

		distance := core.Distance(recv.rocketReceiver.GPS.Latitude, recv.rocketReceiver.GPS.Longitude,
			recv.rocketReceiver.RocketInfo.GPS.Latitude, recv.rocketReceiver.RocketInfo.GPS.Longitude)
		distanceStr := fmt.Sprintf("%0f M", distance)
		fmt.Println(distanceStr)
		recv.distanceLabel.SetText(distanceStr)
		if recv.rocketReceiver.RocketInfo.Recording {
			recv.launchModeStateLabel.SetText("recording")
		} else {
			recv.launchModeStateLabel.SetText("ready")
		}
		if time.Now().Sub(recv.rocketReceiver.LastReceived).Seconds() > 10.0 {
			recv.launchModeStateLabel.SetText("not communicating")
		}
		time.Sleep(time.Second) // Would be nicer if we waited on a channel from RocketInfo
	}
}

func updateLabelWithGPS(label *widget.Label, gps nmea.GGA, altitude float64) {
	labelStr := ""
	if gps.FixQuality == nmea.Invalid {
		labelStr = "Invalid GPS"
	} else {
		labelStr = fmt.Sprintf("Lat: %09.5f Lon:%09.5f Alt: %.2fM", gps.Latitude, gps.Longitude, altitude)
	}
	fmt.Println(labelStr)
	label.SetText(labelStr)
}
