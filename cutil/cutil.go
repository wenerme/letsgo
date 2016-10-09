package cutil

//#include "./cutil.h"
import "C"
import "unsafe"

// (**char,n) -> [][]byte
func StrArrToByteSlice(arr **C.char, c int) [][]byte {
    slices := make([][]byte, c)
    for i := 0; i < c; i ++ {
        s := C.StrArrAt(arr, C.int(i))
        l := C.int(C.strlen(s))
        //func C.GoBytes(cArray unsafe.Pointer, length C.int) []byte
        slices[i] = C.GoBytes(unsafe.Pointer(s), l)
    }
    return slices
}
// (**char,n) -> []string
func StrArrToStringSlice(strArr unsafe.Pointer, c int) []string {
    slices := make([]string, c)
    for i := 0; i < c; i ++ {
        slices[i] = string(C.GoString(C.StrArrAt((**C.char)(strArr), C.int(i))))
    }
    return slices
}
