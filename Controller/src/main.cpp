#include <Arduino.h>
#include "Bitmaps/BMPHeaders.h"
#include "Structs_and_Definitions/Commands.h"
#include "Structs_and_Definitions/Configuration.h"
#include "Structs_and_Definitions/Structs.h"
#include "WiFiLogger/Logger.hpp"
#include <SdFat.h>
#include <stdint.h>
#include <stdlib.h>
#include <GyverOLED.h>
#include <Wire.h>

#define LED_ON() digitalWrite(LED, HIGH)
#define LED_OFF() digitalWrite(LED, LOW)
#define FLASH_LED(x) LED_ON(); delay(x); LED_OFF(); delay(x);

SdFat SDCard;
SdFile ImageFile;
String FileName = "DefaultName.txt";
bool AutoEncode = false;
bool FileOpened = false;
GyverOLED<SSH1106_128x64> Display;
MetaData ImageMetaData = {.res = RESOLUTION::QQVGA, .col = IMAGE::GREYSCALED};

// & Function Declaration
inline byte WaitForCommand(enum COMMAND command) __attribute__((__always_inline__));

// & Function Definitions
byte WaitForCommand(enum COMMAND command) {
  while (true) {
    if (Serial.available() >= 2) {
      enum COMMAND CommandType = static_cast<enum COMMAND>(Serial.read());
      byte DataByte = Serial.read();

      if (CommandType == command) {
        return DataByte;
      }
    }
  }
}

// & Main Program
void setup() {
  Serial.begin(BAUDRATE);
  delay(10);

  // ~ Wire.begin(2, 4);
  // ~ Display.init();

  pinMode(LED, OUTPUT);
  if (!SDCard.begin(SD_CHIP_SELECT, SD_SCK_MHZ(12))) {
    LED_ON();
    while (true) yield();
  }
}

void loop() {
  if (Serial.available() >= 2) {
    enum COMMAND CommandType = static_cast<enum COMMAND>(Serial.read());
    byte DataByte = Serial.read();

    switch (CommandType) {
      case COMMAND_SET_FILE_NAME: {
        FileName = Serial.readStringUntil('\n');
        break;
      }

      case COMMAND_OPEN_FILE: {
        if (!FileOpened) {
          ImageFile.open(FileName.c_str(), O_CREAT | O_WRITE);
          FileOpened = true;
        }

        break;
      }

      case COMMAND_CLOSE_FILE: {
        if (FileOpened) {
          ImageFile.close();
          FileOpened = false; 
        }
        
        break;
      }

      case COMMAND_WRITE_BYTE: {
        if (FileOpened) {
          if (AutoEncode) {
            switch (ImageMetaData.col) {
              case GREYSCALED: {
                byte pixel_bytes[3] = {DataByte, DataByte, DataByte};
                ImageFile.write( pixel_bytes, 3 );
                break;
              }

              case RGB: {
                byte HB = DataByte;
                byte LB = WaitForCommand(COMMAND_WRITE_BYTE);
                uint16_t Pixel = HB << 8 | LB;

                byte r8 = Pixel >> 11 & 0x1F * 255 / 31;
                byte g8 = Pixel >> 5 & 0x3F * 255 / 62;
                byte b8 = Pixel & 0x1F * 255 / 31;

                byte pixel_bytes[3] = {r8, g8, b8};
                ImageFile.write( pixel_bytes, 3 );

                break;
              }

              case BAYER: {
                // ! Not Implemented
                break;
              }

              default: break;
            }
          } else { ImageFile.write(DataByte); }
        }
        break;
      }

      case COMMAND_AUTO_ENCODE: {
        AutoEncode = !AutoEncode;
        break;
      }

      case COMMAND_METADATA_BYTE: {
        ImageMetaData.res = static_cast<enum RESOLUTION>((DataByte & 0b00000111) / 2);
        ImageMetaData.col = static_cast<enum IMAGE>((DataByte & 0b00111000) / 2);

        if (FileOpened) {
          switch (ImageMetaData.res) {
            case QQVGA: {
              ImageFile.write(BmpHeader160x120, BMP_HEADER_SIZE);
              break;
            }

            case QVGA: {
              ImageFile.write(BmpHeader320x240, BMP_HEADER_SIZE);
              break;
            }

            case VGA: {
              ImageFile.write(BmpHeader640x480, BMP_HEADER_SIZE);
              break;
            }

            default: break;
          }
        }
        break;
      }

      case COMMAND_LED_ON: {
        digitalWrite(LED, HIGH);
        break;
      }

      case COMMAND_LED_OFF: {
        digitalWrite(LED, LOW);
        break;
      }

      default: break;
    }
  }
}
