package main

import (
	"log"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, func() {
		log.Println("Exiting application")
	})
}
