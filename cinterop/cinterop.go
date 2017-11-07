package cinterop

/*
#include <string.h>
#include <stdlib.h>
#include <stdint.h>

int cinterop_int_array_at(int* ptr, int i)
{
	return ptr[i];
}

char *cinterop_str_array_at(char **ptr, int i) {
  return ptr[i];
}
*/
import "C"
import "unsafe"

var _ = C.int(0)

// (*int,n) -> []int
func GoIntSlice(ptr unsafe.Pointer, length int) []int {
	ret := make([]int, length)
	for i := 0; i < length; i++ {
		ret[i] = int(C.cinterop_int_array_at((*C.int)(ptr), C.int(i)))
	}
	return ret
}

// (**char,n) -> []string
func GoStringSlice(ptr unsafe.Pointer, length int) []string {
	ret := make([]string, length)
	for i := 0; i < length; i++ {
		ret[i] = string(C.GoString(C.cinterop_str_array_at((**C.char)(ptr), C.int(i))))
	}
	return ret
}

func GoStringSliceN(ptr unsafe.Pointer, lengths []int) []string {
	ret := make([]string, len(lengths))
	for i, v := range lengths {
		ret[i] = string(C.GoStringN(C.cinterop_str_array_at((**C.char)(ptr), C.int(i)), C.int(v)))
	}
	return ret
}
