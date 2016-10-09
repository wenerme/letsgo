package cutil

//#include "./cutil.c"
import "C"
import "unsafe"

func StrArrToByteSlice(arr **C.char, c int) [][]byte {
    slices := make([][]byte, c)
    for i := 0; i < c; i ++ {
        s := C.At(arr, C.int(i))
        l := C.int(C.strlen(s))
        //func C.GoBytes(cArray unsafe.Pointer, length C.int) []byte
        slices[i] = C.GoBytes(unsafe.Pointer(s), l)
    }
    return slices
}
func StrArrToStringSlice(arr **C.char, c int) []string {
    slices := make([]string, c)
    for i := 0; i < c; i ++ {
        s := C.StrArrAt(arr, C.int(i))
        l := C.int(C.strlen(s))
        slices[i] = string(C.GoBytes(unsafe.Pointer(s), l))
    }
    return slices
}
