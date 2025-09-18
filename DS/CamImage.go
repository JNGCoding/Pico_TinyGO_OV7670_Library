package DataStructures

import (
	Camera7670 "PICO_OV7670/Camera"
	"fmt"
	"machine"
	"time"

	"tinygo.org/x/drivers/sdcard"
)

// & Special Settings
const (
	OVERRIDE_SIZE = false
)

/*
 * @brief = Is a type of data structure that is made to store image data.
 * @element ImageType = Stores the format of image.
 * @element Resolution = Stores the size of image.
 * @element ImageData = stores the actual raw data of image.
 */
type CameraImage struct {
	ImageType  Camera7670.IMAGE
	Resolution Camera7670.RESOLUTION
	ImageData  []byte
}

/*
 * @brief = Similar to CameraImage but uses a fixed size queue in an attempt that maybe I am able to use the second core to format the Image data to foreg. PNG or JPEG format
 * @element ImageType = Stores the format of image.
 * @element Resolution = Stores the size of image.
 * @element ImageData = stores the actual raw data of image.
 */
type QueuedCameraImage struct {
	ImageType  Camera7670.IMAGE
	Resolution Camera7670.RESOLUTION
	ImageData  *Queue[byte]
}

/*
 * @brief = Gets the dimension according to the available Resolution Mode.
 * @param image_res = Image Resolution Mode
 * @return = Tuple of integers representing image size.
 */
func get_dimensions(image_res Camera7670.RESOLUTION) (int, int) {
	var W, H int
	switch image_res {
	case Camera7670.VGA:
		W = 640
		H = 480
		break
	case Camera7670.QVGA:
		W = 320
		H = 240
		break
	case Camera7670.QQVGA:
		W = 160
		H = 120
		break
	default:
		W = 640
		H = 480
		break
	}

	return W, H
}

/*
 * @brief = Returns the Byte needed to complete per pixel according to the type of image OV7670 is configured to return.
 * @param image_type = Type of image for eg. RGB565, Bayer or YUV422.
 * @return = Bytes Per pixel
 */
func get_image_type(image_type Camera7670.IMAGE) int {
	bytes_per_pixel := 1
	switch image_type {
	case Camera7670.GREYSCALED, Camera7670.BAYER:
		bytes_per_pixel = 1
		break
	case Camera7670.RGB:
		bytes_per_pixel = 2
		break
	default:
		bytes_per_pixel = 1
		break
	}

	return bytes_per_pixel
}

/*
* @brief = Creates the CameraImage DataStructure.
* @param image_type = The format of image.
* @param image_res = The resolution of image.
* @return = An instance of CameraImage with all the maths sorted out.
! Handle Error.
*/
func CreateImage(image_type Camera7670.IMAGE, image_res Camera7670.RESOLUTION) (*CameraImage, error) {
	bytes_per_pixel := get_image_type(image_type)
	W, H := get_dimensions(image_res)
	if image_res == Camera7670.VGA && !OVERRIDE_SIZE {
		return nil, fmt.Errorf("Impossible to store VGA Size Image in RAM.")
	}
	Data := make([]uint8, W*H*bytes_per_pixel)
	return &CameraImage{ImageType: image_type, Resolution: image_res, ImageData: Data}, nil
}

/*
* @brief = Reads Data from Camera (OV7670 Object) in its ImageData buffer.
* @param Cam = pointer to the OV7670 Object.
* @param SafeMode = Whether to check for image corruption or not if found a corruption raise an error.
* @return = error if found corruption else nil.
! Handle Error.
*/
func (CamImage *CameraImage) ReadImage(Cam *Camera7670.OV7670, SafeMode bool) error {
	bytesPerPixel := get_image_type(CamImage.ImageType)
	width, height := get_dimensions(CamImage.Resolution)
	if CamImage.Resolution == Camera7670.VGA && !OVERRIDE_SIZE {
		return fmt.Errorf("Impossible to store VGA Size Image in RAM.")
	}
	DataCounter := 0

	Cam.WaitForNewFrame() // Wait for VSync high then low

	for row := 0; row < height; row++ {
		if SafeMode {
			Cam.WaitForHorizontalSyncLow()
			for !Cam.HSync.Get() {
				if Cam.VSync.Get() {
					return fmt.Errorf("Corrupted Image. Image till Height = %d is done", row)
				}
			}
		}

		for column := 0; column < width; column++ {
			Cam.WaitForPixelClockLow()
			CamImage.ImageData[DataCounter] = Cam.ReadPins()
			DataCounter++
			Cam.WaitForPixelClockHigh()

			Cam.WaitForPixelClockLow()
			if bytesPerPixel == 2 {
				CamImage.ImageData[DataCounter] = Cam.ReadPins()
				DataCounter++
			}
			Cam.WaitForPixelClockHigh()
		}
	}

	return nil
}

/*
* @brief = Creates the QueuedCameraImage DataStructure.
* @param image_type = The format of image.
* @param image_res = The resolution of image.
* @return = An instance of QueuedCameraImage with all the maths sorted out.
! Handle Error.
*/
func CreateQueuedImage(image_type Camera7670.IMAGE, image_res Camera7670.RESOLUTION) (*QueuedCameraImage, error) {
	bytes_per_pixel := get_image_type(image_type)
	W, H := get_dimensions(image_res)
	if image_res == Camera7670.VGA && !OVERRIDE_SIZE {
		return nil, fmt.Errorf("Impossible to store VGA Size Image in RAM.")
	}

	Data := NewQueue[byte](W * H * bytes_per_pixel)
	return &QueuedCameraImage{ImageType: image_type, Resolution: image_res, ImageData: Data}, nil
}

/*
* @brief = Reads Data from Camera (OV7670 Object) in its ImageData buffer.
* @param Cam = pointer to the OV7670 Object.
* @param SafeMode = Whether to check for image corruption or not if found a corruption raise an error.
* @return = error if found corruption else nil.
! Handle Error.
*/
func (CamImage *QueuedCameraImage) ReadImage(Cam *Camera7670.OV7670, SafeMode bool) error {
	bytesPerPixel := get_image_type(CamImage.ImageType)
	width, height := get_dimensions(CamImage.Resolution)

	if CamImage.Resolution == Camera7670.VGA && !OVERRIDE_SIZE {
		return fmt.Errorf("Impossible to store VGA Size Image in RAM.")
	}

	Cam.WaitForNewFrame() // Checks for VSync pin to go high then low.
	for row := 0; row < height; row++ {
		if SafeMode {
			Cam.WaitForHorizontalSyncLow()
			for !Cam.HSync.Get() {
				if Cam.VSync.Get() {
					return fmt.Errorf("Corrupted Image. Image till Height = %d is done", row)
				}
			}
		}

		for column := 0; column < width; column++ {
			Cam.WaitForPixelClockLow()
			CamImage.ImageData.Enqueue(Cam.ReadPins())
			Cam.WaitForPixelClockHigh()

			Cam.WaitForPixelClockLow()
			if bytesPerPixel == 2 {
				CamImage.ImageData.Enqueue(Cam.ReadPins())
			}
			Cam.WaitForPixelClockHigh()
		}
	}

	return nil
}

/*
* @brief = Flash an image to the USB Interface of pico.
* @param Cam = A pointer to a OV7670 Driver Object.
* @param ImageType = Stores the format of image.
* @param Resolution = Stores the size of image.
* @param SafeMode = Check for image corruption.
* @return = Error if caught any.
! Handle Error.
*/
func FlashImage(Cam *Camera7670.OV7670, ImageType Camera7670.IMAGE, Resolution Camera7670.RESOLUTION, SafeMode bool) error {
	if Resolution == Camera7670.VGA && !OVERRIDE_SIZE {
		return fmt.Errorf("Impossible to store VGA Size Image in RAM.")
	}

	CamImage, _ := CreateImage(ImageType, Resolution)
	if err := CamImage.ReadImage(Cam, SafeMode); err != nil {
		for _, item := range CamImage.ImageData {
			machine.USBCDC.WriteByte(item)
			time.Sleep(time.Microsecond)
		}

		clear(CamImage.ImageData)
		CamImage = nil

		return nil
	} else {
		return err
	}
}

/*
* @brief = Flash an image to the UART Interface of pico.
* @param Cam = A pointer to a OV7670 Driver Object.
* @param ImageType = Stores the format of image.
* @param Resolution = Stores the size of image.
* @param SafeMode = Check for image corruption.
* @return = Error if caught any.
! Handle Error.
*/
func FlashImageToUART(UART *machine.UART, Cam *Camera7670.OV7670, ImageType Camera7670.IMAGE, Resolution Camera7670.RESOLUTION, SafeMode bool) error {
	bytesPerPixel := get_image_type(ImageType)
	width, height := get_dimensions(Resolution)

	Cam.WaitForNewFrame() // Checks for VSync pin to go high then low.
	for row := 0; row < height; row++ {
		if SafeMode {
			Cam.WaitForHorizontalSyncLow()
			for !Cam.HSync.Get() {
				if Cam.VSync.Get() {
					return fmt.Errorf("Corrupted Image. Image till Height = %d is done", row)
				}
			}
		}

		for column := 0; column < width; column++ {
			Cam.WaitForPixelClockLow()
			UART.WriteByte(Cam.ReadPins())
			Cam.WaitForPixelClockHigh()

			Cam.WaitForPixelClockLow()
			if bytesPerPixel >= 2 {
				UART.WriteByte(Cam.ReadPins())
			}
			Cam.WaitForPixelClockHigh()
		}
	}

	return nil
}

/*
* @brief = Store an image in a SD Card if possible.
! Handle Error.
*/
func StoreImage(Cam *Camera7670.OV7670, SDCard *sdcard.Device, Address int64, Resolution Camera7670.RESOLUTION, ImageType Camera7670.IMAGE, SafeMode bool) error {
	width, height := get_dimensions(Resolution)
	bytesPerPixel := get_image_type(ImageType)
	row_buffer := make([]byte, width*bytesPerPixel)
	DataCounter := 0

	Cam.WaitForNewFrame() // Wait for VSync high then low

	for row := 0; row < height; row++ {
		if SafeMode {
			Cam.WaitForHorizontalSyncLow()
			for !Cam.HSync.Get() {
				if Cam.VSync.Get() {
					return fmt.Errorf("Corrupted Image. Image till Height = %d is done", row)
				}
			}
		}

		for column := 0; column < width; column++ {
			Cam.WaitForPixelClockLow()
			row_buffer[DataCounter] = Cam.ReadPins()
			DataCounter++
			Cam.WaitForPixelClockHigh()

			Cam.WaitForPixelClockLow()
			if bytesPerPixel >= 2 {
				row_buffer[DataCounter] = Cam.ReadPins()
				DataCounter++
			}
			Cam.WaitForPixelClockHigh()
		}

		SDCard.WriteAt(row_buffer, Address+int64(row*width))
		DataCounter = 0
	}

	return nil
}
