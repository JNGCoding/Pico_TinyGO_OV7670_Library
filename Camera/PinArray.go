package Camera7670

import "machine"

/*
~ File Description
^ Implements a Datatype that can manipulate and read 8 Pins at a time.
^ It can read the 8 Pin States and package it into a byte.
^ Also it can write to 8 Pins simultaneously with uint8 or byte as a argument.
*/

type PArray struct {
	Pins []machine.Pin
	Mode machine.PinMode
}

func GeneratePArray(pins []machine.Pin, mode machine.PinMode) *PArray {
	return &PArray{Pins: pins, Mode: mode}
}

func (pa *PArray) Init() {
	for _, item := range pa.Pins {
		item.Configure(machine.PinConfig{Mode: pa.Mode})
	}
}

func (pa *PArray) Read() uint8 {
	var result uint8 = 0
	for index, item := range pa.Pins {
		if item.Get() {
			result |= 1 << index
		}
	}

	return result
}

func (pa *PArray) Write(data uint8) {
	for i := 0; i < 8; i++ {
		if data&(1<<i) != 0 {
			pa.Pins[i].High()
		} else {
			pa.Pins[i].Low()
		}
	}
}
