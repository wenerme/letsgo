#ifndef CGO_CUTIL_H
#define CGO_CUTIL_H

#include <string.h>
#include <stdlib.h>

// 避免在 Go 进行指针运算
char *StrArrAt(char **strArr, int n) {
  return strArr[n];
}

intptr_t PtrToIntptr(void*ptr){
    return (intptr_t)(ptr);
}
#endif

