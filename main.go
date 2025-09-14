//go:build tinygo

package main

import (
	CORE "PICO_OV7670/CoreFiles"
	DataStructures "PICO_OV7670/DS"
	"PICO_OV7670/Modules"
	"fmt"
	"machine"
	"runtime"

	"tinygo.org/x/drivers/hd44780i2c"
	"tinygo.org/x/drivers/sdcard"
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

// ^ SPI Configuration
const (
	SD_SCK = machine.NoPin
	SD_SDI = machine.NoPin
	SD_SDO = machine.NoPin
	SD_CS  = machine.NoPin
)

// ^ IMAGE Configuration
const (
	ImageResolution = Modules.QQVGA
	ImageColorSpace = Modules.GREYSCALED
)

// * Variables
var INBUILT_LED machine.Pin
var Application *CORE.Program
var Display hd44780i2c.Device
var _I2C *machine.I2C
var _SPI *machine.SPI
var Camera *Modules.OV7670
var FrameCounter int
var Image *DataStructures.CameraImage
var SDCard sdcard.Device
var MemoryStatus *runtime.MemStats

func main() {
	Application = CORE.CreateApplication()
	Application.LetSetup(func() {
		FrameCounter = 1
		MemoryStatus = &runtime.MemStats{}

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
		Display = hd44780i2c.New(_I2C, 0x27)
		if err := Display.Configure(hd44780i2c.Config{
			Width:       16,
			Height:      2,
			CursorOn:    false,
			CursorBlink: false,
		}); err != nil {
			Application.Exit(1, fmt.Sprintf("Failed to Configure HD44780I2C. Error = %v\n", err))
		}
		Display.DisplayOn(true)
		Display.BacklightOn(true)
		Display.ClearDisplay()
		Display.Home()

		// ^ Camera
		VSync.Configure(machine.PinConfig{Mode: machine.PinInput})
		HSync.Configure(machine.PinConfig{Mode: machine.PinInput})
		MCLK.Configure(machine.PinConfig{Mode: machine.PinPWM})
		PCLK.Configure(machine.PinConfig{Mode: machine.PinInput})
		DataPins.Init()

		Camera = Modules.CreateOV7670(
			_I2C,
			VSync,
			HSync,
			MCLK,
			PCLK,
			DataPins,
		)
		Camera.Initialize(20_000_000)
		if err := Camera.Configure(ImageColorSpace, ImageResolution); err != nil {
			Application.Exit(1, fmt.Sprintf("Failed to Configure Camera. Error = %v\n", err))
		}
		Camera.SetPCLKSpeed(Modules.PCLK_DIV1) // * Writes at 0x11 Register Changes the speed of PCLK giving more time to PICO to scan a pixel. Currently it is at the highest value of 0x1F but you can decrease it to further speedify things.

		// ^ Camera Image Holder
		Image, _ = DataStructures.CreateImage(ImageColorSpace, ImageResolution)

		// ^ INBUILD LED
		INBUILT_LED = CORE.CreateIOPin(25, machine.PinOutput)

		// ^ SD Card

		/*
			SDCard = sdcard.New(machine.SPI1, SD_SCK, SD_SDO, SD_SDI, SD_CS)
			if err := SDCard.Configure(); err != nil {
			}
			machine.SPI1.Configure(machine.SPIConfig{
				Frequency: 8_000_000,
				SCK:       SD_SCK,
				SDO:       SD_SDO,
				SDI:       SD_SDI,
			})
		*/
	})

	Application.LetLoop(func() {
		Display.Home()
		Display.Print([]byte(fmt.Sprintf("FM: %d", FrameCounter)))

		Display.SetCursor(0, 1)
		runtime.ReadMemStats(MemoryStatus)
		Display.Print([]byte(fmt.Sprintf("ALR: %d", MemoryStatus.Alloc)))

		if MemoryStatus.Alloc > 190000 {
			runtime.GC()
		}

		Image.ReadImage(Camera, false)
		/*
			CORE.Print("IMAGE START!")
			Image.ReadImage(Camera, false)
			for i := 0; i < len(Image.ImageData); i++ {
				CORE.WriteByte(Image.ImageData[i])
				time.Sleep(time.Microsecond)
			}
			CORE.Print("IMAGE END!")
		*/
		FrameCounter++
	})

	Application.Run()

	code, message := Application.ExitInfo()
	CORE.PrintLN(fmt.Sprintf("Exit Code - %d, Exit String - %s", code, message))
}
