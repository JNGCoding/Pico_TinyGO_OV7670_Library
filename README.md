#Pico_TinyGO_OV7670_Library#
A modular TinyGO library for interfacing the OV7670 camera module with the Raspberry Pi Pico. Built for flexibility and performance, this library enables image capture across various resolutions and color formats, making it ideal for embedded vision applications.

##📸 Features
- Supports dynamic resolution configuration (e.g., 160×120 to 640×480 and beyond)
- Enables both grayscale and RGB color space image capture
- Transmits raw image data over Serial for host-side processing
- Designed for extensibility: bitmap conversion, SD card logging, and desktop visualization
- Written in idiomatic TinyGO with modular architecture for easy integration

##🚀 Trial Success
Successfully captured and transmitted a 320×240 grayscale image from the OV7670 to a computer. Post-processing includes bitmap generation to make the image window-accessible.

##🔧 Getting Started
This library is ideal for embedded developers experimenting with image sensors in TinyGO. It serves as a foundation for advanced features like real-time streaming, color image processing, and cross-platform logging.
