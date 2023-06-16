package main

import (
	"fmt"
	"github.com/igeekinc/go-rocket/pkg/core"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("usage: rocket-log-analyzer <log file>")
	}
	logFilePath := os.Args[1]
	ril, err := core.ReadRocketInfoLog(logFilePath)
	if err != nil {
		log.Fatalf("Failed reading log file %s with err %v\n", logFilePath, err)
	}

	fmt.Printf("0: GPS Height %f, Baro height:%f\n", ril.Record(0).GPS.Altitude, ril.Record(0).Altitude)
	fmt.Printf("BaseAGL:%f, OffsetAGL:%f\n", ril.BaseAGL(), ril.BaroHeightAGL(0))
	maxHeight, maxHeightRecord := ril.MaxBaroHeightAGL()
	fmt.Printf("Max Height AGL: %f, %v\n", maxHeight, maxHeightRecord)
}
