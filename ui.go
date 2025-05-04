package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers"

	ui "github.com/itohio/tinygui"
)

func hLine(d drivers.Displayer, x, y, w int16, c color.RGBA) {
	for w > 0 {
		d.SetPixel(x, y, c)
		x++
		w--
	}
}
func vLine(d drivers.Displayer, x, y, h int16, c color.RGBA) {
	for h > 0 {
		d.SetPixel(x, y, c)
		y++
		h--
	}
}

func button(p machine.Pin) time.Duration {
	now := time.Now()
	time.Sleep(time.Millisecond * 10)
	for !p.Get() {
		time.Sleep(time.Millisecond)
		machine.Watchdog.Update()
	}
	return time.Since(now)
}

func runButtons() chan ui.UserCommand {
	command := make(chan ui.UserCommand)

	cmd := func(c ui.UserCommand) {
		select {
		case command <- c:
		default:
		}
	}

	go func() {
		N := 0
		for {
			switch {
			case !btnUp.Get():
				d := button(btnUp)
				if d > time.Second*5 {
					cmd(ui.RESET)
				} else if d > time.Second {
					cmd(ui.LONG_UP)
				} else {
					cmd(ui.PREV)
				}
			case !btnDown.Get():
				d := button(btnDown)
				if d > time.Second {
					cmd(ui.LONG_DOWN)
				} else {
					cmd(ui.NEXT)
				}
			case !btnSelect.Get():
				d := button(btnSelect)
				if d > time.Second {
					cmd(ui.ESC)
				} else {
					cmd(ui.ENTER)
				}
			}
			time.Sleep(time.Millisecond * 10)
			if N == 0 {
				cmd(ui.IDLE)
			}
			N = (N + 1) % 100
		}
	}()

	return command
}

func runUI(cmd chan ui.UserCommand, cfg *Config, w *ui.ContainerBase[ui.Widget]) {
	for {
		select {
		case c := <-cmd:
			machine.Watchdog.Update()
			if w.Interact(c) {
				continue
			}

			switch c {
			case ui.RESET:
				machine.CPUReset()
			case ui.LONG_DOWN:
				saveConfig(*cfg)
			case ui.LONG_UP:
				for _, i := range w.Items {
					if w, ok := i.(*Waterer); ok {
						w.active = !w.active
						w.pump.Low()
					}
				}
			}
		}
	}
}
