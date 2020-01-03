package main

import (
	go_rocket "github.com/igeekinc/go-rocket"
	"log"
)

func main() {
	gr, err := go_rocket.InitGPSReader("/dev/ttyS0", 9600, 8, 1)
	if err != nil {
		log.Fatal(err)
	}
	err = gr.UpdateFromGPSLoop()
	if err != nil {
		log.Fatal(err)
	}
}