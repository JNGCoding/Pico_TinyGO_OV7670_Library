package Modules

/*
~ File Description:
^ Simply simulates the functionality of enums to configure OV7670 with ease.
^ Simulates enums named IMAGE (image format) and RESOLUTION (output resolution).
*/

type IMAGE int

const (
	GREYSCALED = iota
	RGB
	BAYER
)

func (i IMAGE) String() string {
	switch i {
	case GREYSCALED:
		return "GREYSCALED"
	case RGB:
		return "RGB"
	case BAYER:
		return "BAYER"
	}

	return "NOT VALID"
}

type RESOLUTION int

const (
	VGA = iota
	QVGA
	QQVGA
)

func (r RESOLUTION) String() string {
	switch r {
	case VGA:
		return "VGA"
	case QVGA:
		return "QVGA"
	case QQVGA:
		return "QQVGA"
	}

	return "NOT VALID"
}
