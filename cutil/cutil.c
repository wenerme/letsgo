#ifndef _CGO_CUTIL
#define _CGO_CUTIL

#include <string.h>

// 避免在 Go 进行指针运算
char *StrArrAt(char **strArr, int n) {
  return strArr[n];
}
#endif

