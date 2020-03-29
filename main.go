package main

import (
	"github.com/go-vgo/robotgo"
)

func main() {
	//log.Println("Starting server..")
	//server.Start()
	robotgo.ScrollMouse(10, "up")
	robotgo.MouseClick("right", true)
	robotgo.MoveMouseSmooth(100, 200, 1.0, 100.0)
}
