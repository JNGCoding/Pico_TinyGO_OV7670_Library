//go:build tinygo

package main

import (
	Camera7670 "PICO_OV7670/Camera"
	CORE "PICO_OV7670/CoreFiles"
	DataStructures "PICO_OV7670/DS"
	"fmt"
	"machine"
	"runtime"
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

var DataPins *Camera7670.PArray = Camera7670.GeneratePArray(
	[]machine.Pin{
		machine.GPIO6,
		machine.GPIO7,
		machine.GPIO10,
		machine.GPIO11,
		machine.GPIO2,
		machine.GPIO3,
		machine.GPIO21,
		machine.GPIO20,
	},
	machine.PinInput,
)

// ^ I2C Configuration
const (
	I2C_SDA  = machine.I2C0_SDA_PIN
	I2C_SCL  = machine.I2C0_SCL_PIN
	I2C_MODE = machine.I2CModeController
)

// ^ UART Configuration
const (
	UART_TX = machine.UART0_TX_PIN
	UART_RX = machine.UART0_RX_PIN
)

// ^ SPI Configuration
const (
	SD_SCK = machine.SPI0_SCK_PIN
	SD_SDO = machine.SPI0_SDO_PIN
	SD_SDI = machine.SPI0_SDI_PIN
	SD_CS  = machine.GPIO17
)

// ^ IMAGE Configuration
const (
	ImageResolution = Camera7670.QQVGA
	ImageColorSpace = Camera7670.GREYSCALED
	PCLKSPEED = Camera7670.PCLK_DIV3
)

// * Variables
var INBUILT_LED machine.Pin
var Application *CORE.Program
var Display hd44780i2c.Device
var Camera *Camera7670.OV7670
var FrameCounter int
var Image *DataStructures.CameraImage
var MemoryStatus *runtime.MemStats

func main() {
	Application = CORE.CreateApplication()
	Application.LetSetup(func() {
		FrameCounter = 1
		MemoryStatus = &runtime.MemStats{}

		// & Initializing I2C
		if err := machine.I2C0.Configure(machine.I2CConfig{
			SDA:  I2C_SDA,
			SCL:  I2C_SCL,
			Mode: I2C_MODE,
		}); err != nil {
			Application.Exit(1, fmt.Sprintf("Failed to Configure I2C Bus. Error = %v\n", err))
		}

		// & Initializing UART
		if err := machine.UART0.Configure(machine.UARTConfig{
			BaudRate: 2_000_000,
			TX:       machine.UART0_TX_PIN,
			RX:       machine.UART0_RX_PIN,
		}); err != nil {
			Application.Exit(1, fmt.Sprintf("Failed to Configure UART Bus. Error = %v\n", err))
		}

		// & Initializing Drivers
		// ^ Display
		Display = hd44780i2c.New(machine.I2C0, 0x27)
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

		Camera = Camera7670.CreateOV7670(
			machine.I2C0,
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
		Camera.SetPCLKSpeed(PCLKSPEED) // * Writes at 0x11 Register Changes the speed of PCLK giving more time to PICO to scan a pixel. Currently it is at the highest value of 0x1F but you can decrease it to further speedify things.

		// ^ Camera Image Holder
		Image, _ = DataStructures.CreateImage(ImageColorSpace, ImageResolution)

		// ^ INBUILD LED
		INBUILT_LED = CORE.CreateIOPin(25, machine.PinOutput)
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

		machine.USBCDC.Write([]byte("IMAGE START!"))
		for i := 0; i < len(Image.ImageData); i++ {
			machine.USBCDC.WriteByte(Image.ImageData[i])
			time.Sleep(time.Microsecond)
		}
		machine.USBCDC.Write([]byte("IMAGE END!"))

		FrameCounter++
	})

	Application.Run()

	code, message := Application.ExitInfo()
	CORE.PrintLN(fmt.Sprintf("Exit Code - %d, Exit String - %s", code, message))
}
