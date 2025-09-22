#ifndef MY_STRUCTS_H_
#define MY_STRUCTS_H_

#include <stdint.h>

enum RESOLUTION { VGA, QVGA, QQVGA };
enum IMAGE { GREYSCALED, RGB, BAYER };

typedef struct Metadata_type {
  enum RESOLUTION res;
  enum IMAGE col;
} MetaData;

#endif