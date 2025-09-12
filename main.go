//go:build tinygo

package main

import (
	CORE "PICO_OV7670/CoreFiles"
	DataStructures "PICO_OV7670/DS"
	"PICO_OV7670/Modules"
	"fmt"
	"machine"
	"time"

	"tinygo.org/x/drivers/hd44780i2c"
)

// ~ CONSTANTS
// ^ Camera Configuration
const (
	MCLK  machine.Pin = machine.GPIO15
	PCLK  machine.Pin = machine.GPIO14
	VSync machine.Pin = machine.GPIO13
	HSync machine.Pin = machine.GPIO12
)

var DataPins *Modules.PArray = Modules.GeneratePArray(
	[]machine.Pin{
		machine.GPIO6,
		machine.GPIO7,
		machine.GPIO10,
		machine.GPIO11,
		machine.GPIO2,
		machine.GPIO3,
		machine.GPIO19,
		machine.GPIO18,
	},
	machine.PinInput,
)

// * Variables
var INBUILT_LED machine.Pin
var Application *CORE.Program
var Display hd44780i2c.Device
var _I2C *machine.I2C
var Camera *Modules.OV7670
var FrameCounter int
var Image DataStructures.CameraImage

func main() {
	Application = CORE.CreateApplication()
	Application.LetSetup(func() {
		time.Sleep(10 * time.Second)
		FrameCounter = 0

		// & Initializing I2C
		_I2C = machine.I2C0
		if err := _I2C.Configure(machine.I2CConfig{
			SDA:  machine.I2C0_SDA_PIN,
			SCL:  machine.I2C0_SCL_PIN,
			Mode: machine.I2CModeController,
		}); err != nil {
			Application.Exit(1, fmt.Sprintf("Failed to Configure I2C Bus. Error = %v\n", err))
		}

		// & Initializing Drivers
		// ^ Display
		/*
			Display = hd44780i2c.New(_I2C, 0x27)
			if err := Display.Configure(hd44780i2c.Config{
				Width:       16,
				Height:      2,
				CursorOn:    true,
				CursorBlink: true,
			}); err != nil {
				Application.Exit(1, fmt.Sprintf("Failed to Configure HD44780I2C. Error = %v\n", err))
			}
			Display.DisplayOn(true)
			Display.BacklightOn(true)
			Display.ClearDisplay()
			Display.Home()
		*/

		VSync.Configure(machine.PinConfig{Mode: machine.PinInput})
		HSync.Configure(machine.PinConfig{Mode: machine.PinInput})
		MCLK.Configure(machine.PinConfig{Mode: machine.PinPWM})
		PCLK.Configure(machine.PinConfig{Mode: machine.PinInput})

		// ^ Camera
		Camera = Modules.CreateOV7670(
			_I2C,
			VSync,
			HSync,
			MCLK,
			PCLK,
			DataPins,
		)
		Camera.Initialize(8_000_000)
		if err := Camera.Configure(Modules.GREYSCALED, Modules.QVGA); err != nil {
			Application.Exit(1, fmt.Sprintf("Failed to Configure Camera. Error = %v\n", err))
		}
		Camera.Write(0x11, 0x1F)

		// ^ Camera Image Holder
		Image = *DataStructures.CreateImage(Modules.GREYSCALED, Modules.QVGA)

		// ^ INBUILD LED
		INBUILT_LED = CORE.CreateIOPin(25, machine.PinOutput)
	})

	Application.LetLoop(func() {
		if err := Image.ReadImage(Camera, false); err != nil {
			machine.USBCDC.Write([]byte("Frame Capture Failed."))
		} else {
			machine.USBCDC.Write([]byte("IMAGE START!"))
			for i := 0; i < len(Image.ImageData); i++ {
				machine.USBCDC.Write([]byte{Image.ImageData[i]})
				time.Sleep(time.Microsecond)
			}
			machine.USBCDC.Write([]byte("IMAGE END!"))
			FrameCounter++
		}

		time.Sleep(time.Second)
	})

	Application.Run()

	code, message := Application.ExitInfo()
	machine.USBCDC.Write([]byte(fmt.Sprintf("Exit Code - %d, Exit String - %s", code, message)))
}
