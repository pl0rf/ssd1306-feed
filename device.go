package main

import (
	"image"
	"log"
	"time"

	"github.com/fogleman/gg"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/devices/ssd1306"
	"periph.io/x/periph/host"
)

var (
	fontPath         = "./Gameplay.ttf"
	fontSize float64 = 24
)

// NOTE: Due to the limited refresh rate of the SSD1306, it is possible for
// price updates to outpace the screen update speed. Therefore when attempting
// to update the screen, we really only care about the very latest price value.
// When the price updates are coming at a very rapid rate, don't want to be
// backed up waiting on the screen to update. Maybe I should just use a single
// shared string with a mutex? But I don't want the goroutine reading from the
// websocket to block.

func deviceLoop(output <-chan string, exit <-chan struct{}) {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}

	opts := &ssd1306.Opts{H: 64, W: 128}
	dev, err := ssd1306.NewI2C(bus, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Halt()

	ctx := gg.NewContext(dev.Bounds().Dx(), dev.Bounds().Dy())
	ctx.LoadFontFace(fontPath, fontSize)

	last := ""

	for {
		select {
		case msg := <-output:
			avail := len(output)
			if avail > 0 { // drain channel and only keep the most recent msg
				for i := 1; i <= avail; i++ {
					msg = <-output
				}
			}

			if last != msg {
				last = msg
				// Setup Graphics Context
				ctx.SetRGB(0, 0, 0)
				ctx.Clear()
				ctx.SetRGB(1, 1, 1)

				// Draw string on our context img
				ctx.DrawStringAnchored(msg, 64, 32, 0.5, 0.5)

				// Send img to device
				err := dev.Draw(dev.Bounds(), ctx.Image(), image.ZP)
				if err != nil {
					log.Println("dev.Draw err: ", err)
				}
			}
		case <-exit:
			log.Println("exit device_loop")
			return
		}
	}

}

// fake_device_loop is a test func for running this on my machine without
// having an oled plugged in, simulate screen draws.
func fakeDeviceLoop(output <-chan string, exit <-chan struct{}) {
	last := ""
	for {
		select {
		case msg := <-output:
			avail := len(output)
			if avail > 0 {
				// read avail
				// log.Println("dropping: ", avail)
				for i := 1; i <= avail; i++ {
					// log.Println("disregard: ", msg)
					msg = <-output
				}
			}

			out := formatStr(msg)
			if last != out {
				last = out
				// draw to screen
				log.Println("output: ", out)
			}
			time.Sleep(500 * time.Millisecond)

		case <-exit:
			log.Println("exit fake_device_loop")
			return
		}
	}
}

/*
func debounce(interval time.Duration, input chan string, cb func(arg string)) {
	var item string
	timer := time.NewTimer(interval)
	for {
		select {
		case item = <-input:
			timer.Reset(interval)
		case <-timer.C:
			if item != "" {
				cb(item)
			}
		}
	}
}

func fake_device_draw(s string) {
	log.Println("device output: ", s)
}
*/
