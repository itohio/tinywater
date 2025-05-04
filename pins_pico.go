//go:build rp2040

package main

import "machine"

var (
	lightR = machine.GP7
	lightG = machine.GP8
	lightB = machine.GP6
	led    = machine.LED

	btnUp     = machine.GP1
	btnSelect = machine.GP2
	btnDown   = machine.GP3

	pump01 = machine.GP10
	pump02 = machine.GP11
	pump03 = machine.GP12
	pump04 = machine.GP13
)
