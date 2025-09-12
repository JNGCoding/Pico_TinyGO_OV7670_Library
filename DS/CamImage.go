package DataStructures

import (
	"PICO_OV7670/Modules"
	"fmt"
)

/*
 * @brief = Is a type of data structure that is made to store image data.
 * @element ImageType = Stores the format of image.
 * @element Resolution = Stores the size of image.
 * @element ImageData = stores the actual raw data of image.
 */
type CameraImage struct {
	ImageType  Modules.IMAGE
	Resolution Modules.RESOLUTION
	ImageData  []byte
}

/*
 * @brief = Similar to CameraImage but uses a fixed size queue in an attempt that maybe I am able to use the second core to format the Image data to foreg. PNG or JPEG format
 * @element ImageType = Stores the format of image.
 * @element Resolution = Stores the size of image.
 * @element ImageData = stores the actual raw data of image.
 */
type QueuedCameraImage struct {
	ImageType  Modules.IMAGE
	Resolution Modules.RESOLUTION
	ImageData  *Queue[byte]
}

/*
 * @brief = Gets the dimension according to the available Resolution Mode.
 * @param image_res = Image Resolution Mode
 * @return = Tuple of integers representing image size.
 */
func get_dimensions(image_res Modules.RESOLUTION) (int, int) {
	var W, H int
	switch image_res {
	case Modules.VGA:
		W = 640
		H = 480
		break
	case Modules.QVGA:
		W = 320
		H = 240
		break
	case Modules.QQVGA:
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
func get_image_type(image_type Modules.IMAGE) int {
	bytes_per_pixel := 1
	switch image_type {
	case Modules.GREYSCALED, Modules.BAYER:
		bytes_per_pixel = 1
		break
	case Modules.RGB:
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
 * @return = An instance of CameraImage will all the maths sorted out.
 */
func CreateImage(image_type Modules.IMAGE, image_res Modules.RESOLUTION) *CameraImage {
	bytes_per_pixel := get_image_type(image_type)
	W, H := get_dimensions(image_res)
	Data := make([]uint8, W*H*bytes_per_pixel)
	return &CameraImage{ImageType: image_type, Resolution: image_res, ImageData: Data}
}

/*
* @brief = Reads Data from Camera (OV7670 Object) in its ImageData buffer.
* @param Cam = pointer to the OV7670 Object.
* @param SafeMode = Whether to check for image corruption or not if found a corruption raise an error.
* @return = error if found corruption else nil.
! Handle Error.
*/
func (CamImage *CameraImage) ReadImage(Cam *Modules.OV7670, SafeMode bool) error {
	bytesPerPixel := get_image_type(CamImage.ImageType)
	width, height := get_dimensions(CamImage.Resolution)
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
 * @return = An instance of QueuedCameraImage will all the maths sorted out.
 */
func CreateQueuedImage(image_type Modules.IMAGE, image_res Modules.RESOLUTION) *QueuedCameraImage {
	bytes_per_pixel := get_image_type(image_type)
	W, H := get_dimensions(image_res)
	Data := NewQueue[byte](W * H * bytes_per_pixel)
	return &QueuedCameraImage{ImageType: image_type, Resolution: image_res, ImageData: Data}
}

/*
* @brief = Reads Data from Camera (OV7670 Object) in its ImageData buffer.
* @param Cam = pointer to the OV7670 Object.
* @param SafeMode = Whether to check for image corruption or not if found a corruption raise an error.
* @return = error if found corruption else nil.
! Handle Error.
*/
func (CamImage *QueuedCameraImage) ReadImage(Cam *Modules.OV7670, SafeMode bool) error {
	bytes_per_pixel := get_image_type(CamImage.ImageType)
	width, height := get_dimensions(CamImage.Resolution)
	data_counter := 0

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
			// Extract Image Data
			var pixel_byte byte

			Cam.WaitForPixelClockLow()
			pixel_byte = Cam.ReadPins()
			CamImage.ImageData.Enqueue(pixel_byte)
			data_counter++

			Cam.WaitForPixelClockHigh()
			Cam.WaitForPixelClockLow()
			if bytes_per_pixel == 2 {
				pixel_byte = Cam.ReadPins()
				CamImage.ImageData.Enqueue(pixel_byte)
				data_counter++
			}
			Cam.WaitForPixelClockHigh()
		}
	}

	return nil
}
