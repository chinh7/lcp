#include <stdint.h>
#include <stdlib.h>
#include <string.h>

typedef Event;
typedef uint8_t* pointer;
typedef pointer* lparray;

extern Event Say(lparray message);

lparray to_lparray(char s[]) {
  lparray result = (lparray) malloc(2 * sizeof(pointer));
  result[0] = strlen(s);
  result[1] = (pointer)s;
  return result;
}

int say(int i) {
  Say(to_lparray("Checking"));
  return i;
}
