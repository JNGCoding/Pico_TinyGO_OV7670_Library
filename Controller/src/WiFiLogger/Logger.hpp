#ifndef LOGGER_H_
#define LOGGER_H_

#include <WiFi.h>
#include <stdint.h>

#define LOGPASS 0x01
#define LOGFAIL 0x00
#define TIMEOUT 3000

class WiFiLogger {
public:
    WiFiLogger(IPAddress ip, uint port);
    ~WiFiLogger();

    void init();
    int8_t write(uint8_t* bytes, uint32_t size);
    int8_t write_byte(uint8_t byte);
    int8_t print(const char* str);
    int8_t println(const char* str);
    uint32_t read(uint8_t* buffer, uint32_t size);
    int8_t readString(String& str, char delimiter);
    void close();
private:
    WiFiClient client;
    IPAddress ServerIP;
    uint ServerPort;
};

#endif