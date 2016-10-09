#ifndef CGO_CUTIL_H
#define CGO_CUTIL_H

#include <string.h>

// 避免在 Go 进行指针运算
char *StrArrAt(char **strArr, int n) {
  return strArr[n];
}
#endif

