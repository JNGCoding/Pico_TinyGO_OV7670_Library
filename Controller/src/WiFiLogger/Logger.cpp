#include "Logger.hpp"
#include <string.h>

WiFiLogger::WiFiLogger(IPAddress ip, uint port) : ServerIP(ip), ServerPort(port) {}
WiFiLogger::~WiFiLogger() {
    if (this->client.connected()) {
        this->client.flush();
        this->client.stop();
    }
}

void WiFiLogger::init() {
    this->client.connect(this->ServerIP, this->ServerPort);
}

int8_t WiFiLogger::write(uint8_t* bytes, uint32_t size) {
    if (this->client.connected()) {
        this->client.write(bytes, size);
        return LOGPASS;
    }

    return LOGFAIL;
}

int8_t WiFiLogger::write_byte(uint8_t byte) {
    if (this->client.connected()) {
        this->client.write(byte);
        return LOGPASS;
    }

    return LOGFAIL;
}

int8_t WiFiLogger::print(const char* str) {
    if (this->client.connected()) {
        this->client.write( (const uint8_t*) str, strlen(str));
        return LOGPASS;
    }

    return LOGFAIL;
}

int8_t WiFiLogger::println(const char* str) {
    if (this->client.connected()) {
        this->client.write( (const uint8_t*) str, strlen(str));
        this->write_byte('\n');
        return LOGPASS;
    }

    return LOGFAIL;
}

uint32_t WiFiLogger::read(uint8_t* buffer, uint32_t size) {
    if (this->client.connected()) {
        uint32_t byte_counter = 0;
        uint32_t start_time = millis();
        while (millis() - start_time < TIMEOUT && byte_counter < size) {
            if (this->client.available() > 0) {
                buffer[byte_counter++] = this->client.read();
                start_time = millis();
            }
        }

        return byte_counter;
    }

    return LOGFAIL;
}

int8_t WiFiLogger::readString(String& str, char delimiter) {
    if (this->client.connected()) {
        uint32_t start_time = millis();
        while (millis() - start_time < TIMEOUT) {
            if (this->client.available() > 0) {
                char byte = this->client.read();
                if (byte == delimiter) break; else str.concat(byte);
                start_time = millis();
            }
        }

        return LOGPASS;
    }

    return LOGFAIL;
}

void WiFiLogger::close() {
    if (this->client.connected()) {
        this->client.flush();
        this->client.stop();
    }
}