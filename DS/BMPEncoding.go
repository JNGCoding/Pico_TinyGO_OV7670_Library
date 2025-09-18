package DataStructures

import Camera7670 "PICO_OV7670/Camera"

/*
 * @brief = ImageStream is a object created to easily get an bmp encoded image from CameraImage which stores raw image data.
 * @element data = Stores the raw image data.
 * @element current_index = Used to determine the data we want to send.
 * @elements Resolution, ImageType = Used to determine the type of image.
 * @element eof = Determines whether the image has ended or still has data.
 */
type ImageStream struct {
	data          []byte
	current_index int64
	Resolution    Camera7670.RESOLUTION
	ImageType     Camera7670.IMAGE
	eof           bool
}

// * Necessary BMPHeaders
var BmpHeader160x120 = []byte{
	0x42, 0x4D,
	0x36, 0x2C, 0x00, 0x00,
	0x00, 0x00,
	0x00, 0x00,
	0x36, 0x00, 0x00, 0x00,
	0x28, 0x00, 0x00, 0x00,
	0xA0, 0x00, 0x00, 0x00,
	0x78, 0x00, 0x00, 0x00,
	0x01, 0x00,
	0x18, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x2C, 0x00, 0x00,
	0x13, 0x0B, 0x00, 0x00,
	0x13, 0x0B, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
}

var BmpHeader320x240 = []byte{
	0x42, 0x4D,
	0xDE, 0x82, 0x03, 0x00,
	0x00, 0x00,
	0x00, 0x00,
	0x36, 0x00, 0x00, 0x00,
	0x28, 0x00, 0x00, 0x00,
	0x40, 0x01, 0x00, 0x00,
	0xF0, 0x00, 0x00, 0x00,
	0x01, 0x00,
	0x18, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x82, 0x03, 0x00,
	0x13, 0x0B, 0x00, 0x00,
	0x13, 0x0B, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
}

var BmpHeader640x480 = []byte{
	0x42, 0x4D,
	0x36, 0x6C, 0x0B, 0x00,
	0x00, 0x00,
	0x00, 0x00,
	0x36, 0x00, 0x00, 0x00,
	0x28, 0x00, 0x00, 0x00,
	0x80, 0x02, 0x00, 0x00,
	0xE0, 0x01, 0x00, 0x00,
	0x01, 0x00,
	0x18, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x6C, 0x0B, 0x00,
	0x13, 0x0B, 0x00, 0x00,
	0x13, 0x0B, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
}

/*
 * @brief = Creates a ImageStream object.
 * @param __image__ = A pointer to a CameraImage object.
 * @return = Returns a pointer to an instance of ImageStream Object.
 */
func EncodeImage(__image__ *CameraImage) *ImageStream {
	return &ImageStream{data: __image__.ImageData, current_index: 0, Resolution: __image__.Resolution, ImageType: __image__.ImageType, eof: false}
}

/*
 * @brief = Returns a formatted pixel from the ImageStream object for later use.
 * @return = []byte{} object which contains the pixel data in RGB888 Format.
 */
func (stream *ImageStream) GetNextPixel() []byte {
	increment_flag := true
	var result []byte
	if stream.current_index < 54 {
		switch stream.Resolution {
		case Camera7670.VGA:
			result = []byte{BmpHeader640x480[stream.current_index]}
			break
		case Camera7670.QVGA:
			result = []byte{BmpHeader320x240[stream.current_index]}
			break
		case Camera7670.QQVGA:
			result = []byte{BmpHeader160x120[stream.current_index]}
			break
		default:
			increment_flag = false
			break
		}
	} else if stream.current_index >= 54 && stream.current_index < int64(len(stream.data)) {
		switch stream.ImageType {
		case Camera7670.GREYSCALED:
			result = []byte{
				stream.data[stream.current_index-54],
				stream.data[stream.current_index-54],
				stream.data[stream.current_index-54],
			}
			break
		case Camera7670.RGB:
			var item uint16 = uint16(stream.data[stream.current_index+1])<<8 | uint16(stream.data[stream.current_index])
			var r8, g8, b8 uint8

			r8 = byte(item >> 11 & 0x1F)
			g8 = byte(item >> 5 & 0x3F)
			b8 = byte(item & 0x1F)

			result = []byte{
				r8,
				g8,
				b8,
			}
			break
		case Camera7670.BAYER:
			increment_flag = false // ! NOT IMPLEMENTED
			break
		default:
			increment_flag = false
		}
	} else {
		increment_flag = false
		stream.eof = true
	}

	if increment_flag {
		if stream.ImageType == Camera7670.RGB {
			stream.current_index += 2
		} else {
			stream.current_index++
		}
		return result
	}

	return []byte{0x00}
}

// & OneLine Brief = Checks if the image has ended or not.
func (stream *ImageStream) GetEOF() bool {
	return stream.eof
}
