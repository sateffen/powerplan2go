package main

import (
	_ "embed"
	"log"
	"os/exec"
	"strings"
	"syscall"

	"github.com/getlantern/systray"
	"golang.org/x/text/encoding/charmap"
)

//go:embed icon.ico
var iconIco []byte

// Generic handler for clicks on any powerplan. Enables the powerplan and disables
// all others.
func handleClick(menuItems []*systray.MenuItem, index int, name string, guid string) {
	for {
		<-menuItems[index].ClickedCh
		log.Printf("Clicked %s \n", name)

		cmd := exec.Command("powercfg", "-s", guid)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		err := cmd.Run()

		if err != nil {
			log.Fatalf("Could not activate powerschema %s with guid %s \n", name, guid)
			continue
		}

		for _, item := range menuItems {
			item.Uncheck()
		}

		menuItems[index].Check()
	}
}

// The ready handler for setting up the application after systray startup
// finishes. This extracts the available powerplans
func onReady() {
	log.Println("Starting setup")

	systray.SetIcon(iconIco)
	systray.SetTitle("Powerplan2go")
	systray.SetTooltip("Powerplan2go")

	cmd := exec.Command("powercfg", "-l")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()

	if err != nil {
		log.Panicln("Could not retrieve powerlevels")
	}

	outputAsUTF8, err := charmap.CodePage850.NewDecoder().Bytes(output)

	if err != nil {
		log.Panicln("Decoding of output failed")
	}

	allOutputLines := strings.Split(string(outputAsUTF8), "\n")
	relevantOutputLines := allOutputLines[3 : len(allOutputLines)-1]
	menuItems := make([]*systray.MenuItem, len(relevantOutputLines))

	for i, line := range relevantOutputLines {
		line = strings.TrimSpace(line)

		relevantString := strings.TrimSpace(strings.Split(line, ":")[1])
		parts := strings.Split(relevantString, " ")
		guid := parts[0]
		name := strings.ToValidUTF8(parts[2][1:len(parts[2])-1], "")

		menuItems[i] = systray.AddMenuItem(name, guid)

		if len(parts) == 4 {
			menuItems[i].Check()
		}

		go handleClick(menuItems, i, name, guid)
	}

	systray.AddSeparator()
	quitMenuItem := systray.AddMenuItem("Quit", "Quit Powerplan2go")

	go func() {
		<-quitMenuItem.ClickedCh
		log.Println("Quitting app")
		systray.Quit()
	}()
}
