package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	ui "github.com/itohio/tinygui"
	"github.com/itohio/tinygui/widget"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/drivers/ws2812"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

//go:generate tinygo flash -target=pico

var (
	WIDTH  int16 = 120
	HEIGHT int16 = 8
	day    bool
)
var (
	white = color.RGBA{255, 255, 255, 255}
	black = color.RGBA{0, 0, 0, 0}
)

func main() {
	machine.InitADC()

	lightR.Configure(machine.PinConfig{Mode: machine.PinOutput})
	lightG.Configure(machine.PinConfig{Mode: machine.PinOutput})
	lightB.Configure(machine.PinConfig{Mode: machine.PinOutput})

	btnUp.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	btnDown.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	btnSelect.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.NewWS2812(led)
	ws.WriteColors([]color.RGBA{{R: 255}})

	machine.I2C0.Configure(machine.I2CConfig{Frequency: 400 * machine.KHz})
	// machine.LED.High()

	// the delay is needed for display start from a cold reboot, not sure why
	time.Sleep(time.Second)

	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{Width: 128, Height: 64, Address: 0x3C, VccState: ssd1306.SWITCHCAPVCC})
	display.ClearDisplay()

	cfg, err := loadConfig()
	if err != nil {
		tinyfont.WriteLine(&display, &freemono.Regular9pt7b, 0, 11, err.Error(), white)
		err = saveConfig(cfg)
		if err != nil {
			tinyfont.WriteLine(&display, &freemono.Regular9pt7b, 0, 22, err.Error(), white)
		}
		display.Display()
		time.Sleep(time.Minute)
	}

	var dashboard *ui.ContainerBase[ui.Widget]
	dashboard = ui.NewContainer[ui.Widget](
		uint16(WIDTH), 0, ui.LayoutVList(1),
		NewWaterer(0, machine.ADC0, pump01, &cfg),
		NewWaterer(1, machine.ADC1, pump02, &cfg),
		NewWaterer(2, machine.ADC2, pump03, &cfg),
		NewWaterer(3, machine.ADC3, pump04, &cfg),
		widget.NewLabel(uint16(WIDTH), 10, &tinyfont.TomThumb, func() string {
			active := "o"
			for _, i := range dashboard.Items {
				if w, ok := i.(*Waterer); ok {
					if w.active {
						active = "#"
					}
					break
				}
			}

			param := ""
			for _, i := range dashboard.Items {
				if w, ok := i.(*Waterer); ok {
					if w.Selected() {
						param = fmt.Sprintf(" [%.03f %.03f %.03f]", w.cfg.LowThreshold[w.index], w.lastVal, w.cfg.HighThreshold[w.index])
						break
					}
				}
			}

			return fmt.Sprintf("%s%s", active, param)
		}, white),
	)
	dW, dH := dashboard.Size()
	ctx := ui.NewRandomContext(&display, time.Second*1, dW, dH)

	// Watering logic
	for _, i := range dashboard.Items {
		if w, ok := i.(*Waterer); ok {
			go w.run()
		}
	}

	machine.Watchdog.Configure(machine.WatchdogConfig{
		TimeoutMillis: 3000,
	})
	machine.Watchdog.Start()

	go runDayloop()
	cmd := runButtons()
	go runUI(cmd, &cfg, dashboard)

	// Drawing and moving display around
	ticker := time.NewTicker(time.Millisecond * 50)
	for range ticker.C {
		ws.WriteColors([]color.RGBA{{G: 255}})

		display.ClearBuffer()
		dashboard.Draw(&ctx)
		display.Display()
		ws.WriteColors([]color.RGBA{{G: 0}})
	}
}

func runDayloop() {
	// Illumination
	ticker := time.NewTicker(time.Hour * 24)
	for {
		lightR.High()
		lightG.High()
		lightB.High()
		day = true
		time.Sleep(time.Hour * 9)
		day = false
		lightR.Low()
		lightG.Low()
		lightB.Low()
		<-ticker.C
	}
}
