package SDController

import (
	Camera7670 "PICO_OV7670/Camera"
	"machine"
	"time"
)

// COMMANDS
const (
	COMMAND_SET_FILE_NAME = iota
	COMMAND_OPEN_FILE
	COMMAND_CLOSE_FILE
	COMMAND_WRITE_BYTE
	COMMAND_AUTO_ENCODE
	COMMAND_METADATA_BYTE
	COMMAND_LED_ON
	COMMAND_LED_OFF
)

type SDCard struct {
	CommunicationLine *machine.UART
}

func (SD *SDCard) CreateFile(name string) {
	SD.CommunicationLine.Write([]byte{COMMAND_SET_FILE_NAME, 0x00})
	SD.CommunicationLine.Write([]byte(name + "\n"))
	SD.CommunicationLine.Write([]byte{COMMAND_OPEN_FILE, 0x00})
}

func (SD *SDCard) CloseFile() {
	SD.CommunicationLine.Write([]byte{COMMAND_CLOSE_FILE, 0x00})
}

func (SD *SDCard) Write(data []byte, ChunkSize uint32, DelayLength time.Duration) {
	for i := 0; i < len(data); i++ {
		if i%int(ChunkSize) == 0 {
			time.Sleep(DelayLength)
		}

		SD.CommunicationLine.Write([]byte{COMMAND_WRITE_BYTE, data[i]})
	}
}

func (SD *SDCard) WriteByte(data byte) {
	SD.CommunicationLine.Write([]byte{COMMAND_WRITE_BYTE, data})
}

func (SD *SDCard) TurnOnLED() {
	SD.CommunicationLine.Write([]byte{COMMAND_LED_ON, 0x00})
}

func (SD *SDCard) TurnOffLED() {
	SD.CommunicationLine.Write([]byte{COMMAND_LED_OFF, 0x00})
}

func (SD *SDCard) WriteHeader(Resolution Camera7670.RESOLUTION, ImageType Camera7670.IMAGE) {
	var data byte = 0x00
	data = 1 << (Resolution + 1)
	data = 1 << (ImageType + 4)
	SD.CommunicationLine.Write([]byte{COMMAND_METADATA_BYTE, data})
}

func (SD *SDCard) ToggleAutoEncode() {
	SD.CommunicationLine.Write([]byte{COMMAND_AUTO_ENCODE, 0x00})
}
