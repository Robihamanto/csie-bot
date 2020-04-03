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
	robotgo.ShowAlert("Turn on the server", "Have you complete all setup? State? Mock? Iqomah time?")
	log.Println("Starting server..")
	server.Start(3, 2)
	// x, y := robotgo.GetMousePos()
	// log.Println("X:", x)
	// log.Println("Y:", y)
}
