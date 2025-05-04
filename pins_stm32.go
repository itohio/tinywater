//go:build stm32

package main

import "machine"

// pins not tested
var (
	lightR = machine.PB8
	lightG = machine.PB7
	lightB = machine.PB6
	led    = machine.LED

	btnUp     = machine.PB4
	btnSelect = machine.PB5
	btnDown   = machine.PB6

	pump01 = machine.PB10
	pump02 = machine.PB11
	pump03 = machine.PB12
	pump04 = machine.PB13
)
