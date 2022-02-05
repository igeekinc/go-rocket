package ground

import (
	"fmt"
	"github.com/adrianmo/go-nmea"
	"github.com/igeekinc/go-rocket/pkg/core"
	"html/template"

	"net/http"
	"path"
	"time"
)

/*
This is the HTTP server that runs on the ground.  It is used to show the position of the rocket
 */

type GroundHTTPServer struct {
	httpRoot string
	httpPort int
	receiver * RocketReceiver
	previousPositions []nmea.RMC
}

// NewGroundHTTPServer creates and run the HTTP server.  It does not return unless the HTTP server has an error
func NewGroundHTTPServer(httpRoot string, httpPort int, receiver * RocketReceiver) *GroundHTTPServer {
	ghs := GroundHTTPServer{
		httpRoot:   httpRoot,
		httpPort:   httpPort,
		receiver: receiver,
	}
	http.HandleFunc("/", ghs.indexPage)

	requestMap := map[string]func(writer http.ResponseWriter, request *http.Request){
		"video":func(writer http.ResponseWriter, request *http.Request) {
			ghs.video(writer, request)
		},
	}
	for request, function := range requestMap {
		http.HandleFunc("/api/"+request, function)
	}

	return &ghs
}

func (recv * GroundHTTPServer) UpdateGPSLoop() {
	for (true) {
		curPos := recv.receiver.rocketInfo.GPS
		appendPos := true
		if len(recv.previousPositions) > 0 {
			ourPos := recv.previousPositions[len(recv.previousPositions) - 1]
			distance := core.Distance(ourPos.Latitude, ourPos.Longitude, curPos.Latitude, curPos.Longitude)
			if distance < 10.0 {
				appendPos = false	// We don't want to keep the little movements
			}
		}
		if appendPos {
			recv.previousPositions = append(recv.previousPositions, curPos)
			if len(recv.previousPositions) > 30 {
				recv.previousPositions = recv.previousPositions[1:]
			}
		}
		time.Sleep(30 * time.Second)
	}
}

// Starts the HTTP server - will not return unless there is a major error
func (recv * GroundHTTPServer) Serve() error {
	go recv.UpdateGPSLoop()
	return http.ListenAndServe(fmt.Sprintf(":%d", recv.httpPort), nil)
}

type pageInputs struct {
	CurrentLat, CurrentLong float64
	PreviousPositions []nmea.RMC
}

func (recv * GroundHTTPServer) indexPage(writer http.ResponseWriter, request *http.Request) {
	template, err := template.ParseFiles(path.Join(recv.httpRoot, "index.html"))
	if err != nil {

	}
	currentPos := recv.previousPositions[len(recv.previousPositions) - 1]
	pi := pageInputs {
		CurrentLat: currentPos.Latitude,
		CurrentLong: currentPos.Longitude,
		PreviousPositions: recv.previousPositions,
	}

	template.Execute(writer, &pi)
}

func (recv * GroundHTTPServer) video(writer http.ResponseWriter, request *http.Request) {
	recv.receiver.serialPort.Write([]byte("V\n"))
}