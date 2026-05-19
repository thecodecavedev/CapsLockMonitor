package main

import (
	_ "embed"
	"log"
	"time"

	"github.com/getlantern/systray"
	evdev "github.com/gvalkov/golang-evdev"
)

//go:embed assets/capslock.png
var iconOn []byte

//go:embed assets/capslock_disabled.png
var iconOff []byte

var iconUpdateCh = make(chan bool)
var capLockOn = false

func main() {
	go listenForCapsLock()
	systray.Run(onReady, onExit)

}

func listenForCapsLock() {

	dev, err := evdev.Open("/dev/input/event7")

	if err != nil {
		log.Fatal(err)
	}

	for {
		events, err := dev.Read()

		if err != nil {
			log.Fatal(err)
		}

		for _, ev := range events {
			if ev.Type == evdev.EV_KEY && ev.Code == evdev.KEY_CAPSLOCK && ev.Value == 1 {
				capLockOn = !capLockOn
				iconUpdateCh <- capLockOn
			}
		}
		time.Sleep(10 * time.Millisecond)

	}

}

func onReady() {
	systray.SetIcon(iconOff)

	go func() {
		for capsOn := range iconUpdateCh {
			if capsOn {
				systray.SetIcon(iconOn)
			} else {
				systray.SetIcon(iconOff)

			}
		}
	}()

	systray.SetTitle("Hello Status Bar")
	systray.SetTooltip("Headlines")
}

// func getIcon(s string) []byte {
// 	b, err := os.ReadFile(s)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return b
// }

func onExit() {

}
