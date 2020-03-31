package main

import (
	"log"

	"github.com/Robihamanto/csie-bot/server"
	"github.com/go-vgo/robotgo"
)

func main() {
	// Dont forget to
	// 1. Send state from current time and state
	// Uncomment pMock
	// Change iqomah time
	log.Println("Starting server..")
	server.Start(1)
	x, y := robotgo.GetMousePos()
	log.Println("X:", x)
	log.Println("Y:", y)
}
