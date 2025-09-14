package Modules

/*
~ File Description:
^ Simply simulates the functionality of enums to configure OV7670 with ease.
^ Simulates enums named IMAGE (image format) and RESOLUTION (output resolution).
*/

type PCLK_DIVIDER int

const (
	PCLK_DIV0 = 0x01
	PCLK_DIV1 = 0x04
	PCLK_DIV2 = 0x08
	PCLK_DIV3 = 0x10
	PCLK_DIV4 = 0x1F
)

func (i PCLK_DIVIDER) String() string {
	switch i {
	case PCLK_DIV0:
		return "PCLKDIV0"
	case PCLK_DIV1:
		return "PCLKDIV1"
	case PCLK_DIV2:
		return "PCLKDIV2"
	case PCLK_DIV3:
		return "PCLKDIV3"
	case PCLK_DIV4:
		return "PCLKDIV4"
	}

	return "NOT VALID"
}

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
