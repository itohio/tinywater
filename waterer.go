package main

import (
	"machine"
	"time"

	ui "github.com/itohio/tinygui"
	"tinygo.org/x/drivers"
)

type Waterer struct {
	ui.WidgetBase
	cfg       *Config
	active    bool
	index     byte
	adc       machine.ADC
	pump      machine.Pin
	lastVal   float32
	setHighTh bool
}

func NewWaterer(n byte, adc machine.Pin, pump machine.Pin, cfg *Config) *Waterer {
	ret := &Waterer{
		WidgetBase: ui.NewWidgetBase(uint16(WIDTH), 8),
		cfg:        cfg,
		active:     false,
		index:      n,
		adc:        machine.ADC{Pin: adc},
		pump:       pump,
	}

	pump.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ret.adc.Configure(machine.ADCConfig{})
	return ret
}

func (w *Waterer) run() {
	ticker := time.NewTicker(time.Millisecond * 100)
	for range ticker.C {
		w.read()

		if !w.active {
			continue
		}
		if !w.cfg.IsOn(w.index) {
			continue
		}

		// Don't pump at night!
		if !day {
			w.pump.Low()
			continue
		}

		//
		th := w.cfg.High(w.index)
		if w.pump.Get() {
			th = w.cfg.Low(w.index)
		}

		w.pump.Set(w.lastVal > th && day)
	}
}

func (w *Waterer) read() {
	time.Sleep(time.Microsecond)
	val := float32(w.adc.Get()) / 65535.0
	w.lastVal = (w.lastVal*9 + val) / 10
}

func (w *Waterer) Draw(ctx ui.Context) {
	x, y := ctx.DisplayPos()
	w.display(ctx.D(), x, y)
}

func (w *Waterer) SetSelected(s bool) {
	if s {
		w.setHighTh = false
	}
	w.WidgetBase.SetSelected(s)
}

func (w *Waterer) display(d drivers.Displayer, x, y int16) {
	y += 4
	ui.VLine(d, x, y+1, HEIGHT-2, white)
	ui.VLine(d, x+WIDTH-1, y+1, HEIGHT-2, white)

	if w.Selected() {
		ui.VLine(d, x+1, y, HEIGHT, white)
		ui.VLine(d, x+WIDTH-2, y, HEIGHT, white)
	}

	th := int16(w.cfg.Low(w.index) * float32(WIDTH))
	if w.Selected() && !w.setHighTh {
		ui.VLine(d, x+th, y, HEIGHT, white)
	} else {
		ui.VLine(d, x+th, y+4, 2, white)
	}

	th = int16(w.cfg.High(w.index) * float32(WIDTH))
	if w.Selected() && w.setHighTh {
		ui.VLine(d, x+th, y, HEIGHT, white)
	} else {
		ui.VLine(d, x+th, y+2, 2, white)
	}

	val := w.lastVal
	ui.HLine(d, x, y+4, int16(val*float32(WIDTH)), white)
}

func (w *Waterer) Interact(cmd ui.UserCommand) bool {
	switch cmd {
	case ui.NEXT:
		if w.setHighTh {
			w.cfg.HighThreshold[w.index] += 0.01
		} else {
			w.cfg.LowThreshold[w.index] += 0.01
		}
	case ui.PREV:
		if w.setHighTh {
			w.cfg.HighThreshold[w.index] -= 0.01
		} else {
			w.cfg.LowThreshold[w.index] -= 0.01
		}
	case ui.ENTER:
		if w.setHighTh {
			w.SetSelected(false)
			return true
		}
		w.setHighTh = true
	}

	return w.WidgetBase.Interact(cmd)
}
