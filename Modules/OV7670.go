package Modules

import (
	"fmt"
	"machine"
	"time"
)

// & ImageType and ImageResolution for usage across multiple files.
var CurrentImageType IMAGE
var CurrentResolution RESOLUTION

// & Default address of OV7670 I2C Interface.
const DEFAULT_OV7670_ADDRESS uint8 = 0x21

/*
 * @brief = First Ever TinyGO OV7670 Driver.
 * @element Address = Address of the OV7670 Object.
 * @element I2C_Bus = Pointer to the I2C Driver of TinyGO.
 * @elements VSync, HSync, PCLK = IO Pins to track the image data given by OV7670.
 * @element MCLK = PWM Enabled Pin which is used to generate a clock for the image sensor.
 * @element DataPins = pointer of PArray made to read data from the 8 Data pins with ease.
 */
type OV7670 struct {
	Address  uint8
	I2C_Bus  *machine.I2C
	VSync    machine.Pin
	HSync    machine.Pin
	MCLK     machine.Pin
	PCLK     machine.Pin
	DataPins *PArray
}

/*
 * @brief = Creates a pointer to an instance OV7670 Driver Object.
 * @params = Too many params, not explaining they just go to their corresponding args in the struct construction.
 * @return = pointer of an OV7670 Driver Object.
 */
func CreateOV7670(bus_i2c *machine.I2C, vsync, hsync, mclk, pclk machine.Pin, data_pins *PArray) *OV7670 {
	return &OV7670{Address: DEFAULT_OV7670_ADDRESS, I2C_Bus: bus_i2c, VSync: vsync, HSync: hsync, MCLK: mclk, PCLK: pclk, DataPins: data_pins}
}

/*
* @brief = Generates the Clock Signal and resets the OV7670 Camera.
* @return = returns an error if found any.
! Handle Error.
*/
func (Cam *OV7670) Initialize(_freq uint64) error {
	// Producing MCLK.
	slice, _ := machine.PWMPeripheral(Cam.MCLK)
	pwm := machine.PWM0
	switch slice {
	case 0:
		pwm = machine.PWM0
		break
	case 1:
		pwm = machine.PWM1
		break
	case 2:
		pwm = machine.PWM2
		break
	case 3:
		pwm = machine.PWM3
		break
	case 4:
		pwm = machine.PWM4
		break
	case 5:
		pwm = machine.PWM5
		break
	case 6:
		pwm = machine.PWM6
		break
	case 7:
		pwm = machine.PWM7
		break
	default:
		return fmt.Errorf("Invalid PWM Pin - Failed to create MCLK Signal.")
	}

	var freq uint64 = _freq
	var period = uint64(1e9 / freq)
	if err := pwm.Configure(machine.PWMConfig{Period: period}); err != nil {
		return err
	}

	if ch, err := pwm.Channel(Cam.MCLK); err != nil {
		return err
	} else {
		pwm.Set(ch, pwm.Top()/2)
	}

	// Writing Start-up registers.
	Cam.Reset()
	Cam.Write(0x3A, 0x04)
	Cam.Write(0x13, 0xC0)
	Cam.Write(0x00, 0x00)
	Cam.Write(0x10, 0x00)
	Cam.Write(0x0D, 0x40)
	Cam.Write(0x14, 0x18)
	Cam.Write(0x24, 0x95)
	Cam.Write(0x25, 0x33)
	Cam.Write(0x13, 0xC5)
	Cam.Write(0x6A, 0x40)
	Cam.Write(0x01, 0x40)
	Cam.Write(0x02, 0x60)
	Cam.Write(0x13, 0xC7)
	Cam.Write(0x41, 0x08)
	Cam.Write(0x15, 0x20)

	return nil
}

/*
* @brief = Writes a value to a selected register of OV7670
* @param reg = The Desired Register.
* @param val = The Desired value we want to put to the register.
* @return = Returns I2C Error if any.
! Handle Error.
*/
func (Cam *OV7670) Write(reg, val uint8) error {
	err := Cam.I2C_Bus.WriteRegister(Cam.Address, reg, []byte{val})
	time.Sleep(time.Millisecond)
	return err
}

/*
* @brief = Reads a value to a selected register of OV7670.
* @param reg = The Desired Register.
* @return = Returns the data from the register and I2C Error if any.
! Handle Error.
*/
func (Cam *OV7670) Read(reg uint8) (uint8, error) {
	var result []uint8 = make([]uint8, 1)
	err := Cam.I2C_Bus.ReadRegister(Cam.Address, reg, result)
	return result[0], err
}

/*
 * @brief = Sets the image output resolution of OV7670
 * @param res = The desired resolution.
 */
func (Cam *OV7670) set_resolution(res RESOLUTION) {
	CurrentResolution = res
	switch res.String() {
	case "VGA":
		Cam.Write(0x0C, 0x00)
		Cam.Write(0x32, 0xF6)
		Cam.Write(0x17, 0x13)
		Cam.Write(0x18, 0x01)
		Cam.Write(0x19, 0x02)
		Cam.Write(0x1A, 0x7A)
		Cam.Write(0x03, 0x0A)
		break
	case "QVGA":
		Cam.Write(0x0C, 0x04)
		Cam.Write(0x3E, 0x19)
		Cam.Write(0x72, 0x11)
		Cam.Write(0x73, 0xF1)
		Cam.Write(0x17, 0x16)
		Cam.Write(0x18, 0x04)
		Cam.Write(0x32, 0xA4)
		Cam.Write(0x19, 0x02)
		Cam.Write(0x1A, 0x7A)
		Cam.Write(0x03, 0x0A)
		break
	case "QQVGA":
		Cam.Write(0x19, 0x00)
		Cam.Write(0x1A, 0x7A)
		Cam.Write(0x03, 0x00)
		Cam.Write(0x17, 0x16)
		Cam.Write(0x18, 0x05)
		Cam.Write(0x32, 0x32)
		Cam.Write(0x0C, 0x04)
		Cam.Write(0x3E, 0x1A)
		Cam.Write(0x72, 0x22)
		Cam.Write(0x73, 0xF2)
		break
	default:
		return
	}
}

/*
 * @brief = Sets the image output resolution of OV7670
 * @param res = The desired resolution.
 */
func (Cam *OV7670) set_color(col IMAGE) {
	CurrentImageType = col
	switch col.String() {
	case "GREYSCALED":
		Cam.Write(0x12, 0x00)
		Cam.Write(0x8C, 0x00)
		Cam.Write(0x04, 0x00)
		Cam.Write(0x40, 0xC0)
		Cam.Write(0x14, 0x1A)
		Cam.Write(0x3D, 0x40)
		break
	case "RGB":
		Cam.Write(0x12, 0x04)
		Cam.Write(0x8C, 0x00)
		Cam.Write(0x04, 0x00)
		Cam.Write(0x40, 0xD0)
		Cam.Write(0x14, 0x6A)
		Cam.Write(0x4F, 0xB3)
		Cam.Write(0x50, 0xB3)
		Cam.Write(0x51, 0x00)
		Cam.Write(0x52, 0x3D)
		Cam.Write(0x53, 0xA7)
		Cam.Write(0x54, 0xE4)
		Cam.Write(0x3D, 0x40)
		break
	case "BAYER":
		Cam.Write(0x12, 0x01)
		Cam.Write(0x3D, 0x08)
		Cam.Write(0x41, 0x3D)
		Cam.Write(0x76, 0xE1)
		break
	default:
		return
	}
}

/*
* @brief = Sets the desired Image type and Output Resolution.
* @return = returns an error if the size or image type are not valid.
! Handle Error.
*/
func (Cam *OV7670) Configure(col IMAGE, res RESOLUTION) error {
	if col.String() == "NOT VALID" || res.String() == "NOT VALID" {
		return fmt.Errorf("INVALID COLOR OR RESOLUTION, RESOLUTION: %s | IMAGE: %s\n", res.String(), col.String())
	}

	Cam.set_color(col)
	Cam.set_resolution(res)

	return nil
}

/*
 * @brief = Just returns the PArray Read().
 * @return = returns a byte made from the 8 Pin States.
 */
func (Cam *OV7670) ReadPins() uint8 {
	return Cam.DataPins.Read()
}

/*
& OneLine Brief = Simply Halts the processor till VSync completes its cycle.
*/
func (Cam *OV7670) WaitForNewFrame() {
	for Cam.VSync.Get() {
	}
	for !Cam.VSync.Get() {
	}
}

/*
& OneLine Brief = Simply Halts the processor till PCLK becomes High.
*/
func (Cam *OV7670) WaitForPixelClockLow() {
	for !Cam.PCLK.Get() {
	}
}

/*
& OneLine Brief = Simply Halts the processor till PCLK becomes Low.
*/
func (Cam *OV7670) WaitForPixelClockHigh() {
	for Cam.PCLK.Get() {
	}
}

/*
& OneLine Brief = Simply Halts the processor till HSync becomes High.
*/
func (Cam *OV7670) WaitForHorizontalSyncHigh() {
	for Cam.HSync.Get() {
	}
}

/*
& OneLine Brief = Simply Halts the processor till HSync becomes Low.
*/
func (Cam *OV7670) WaitForHorizontalSyncLow() {
	for !Cam.HSync.Get() {
	}
}

/*
& OneLine Brief = Resets all the registers in OV7670.
*/
func (Cam *OV7670) Reset() {
	Cam.Write(0x12, 0x80)
	time.Sleep(100 * time.Millisecond)
}
