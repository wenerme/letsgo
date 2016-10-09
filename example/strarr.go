package main

/*
char *strs[] = {
    "Hello",
    "There",
    "中文"
};
char **GetStrArr() {
  return strs;
}
 */
import "C"
import (
    "github.com/wenerme/letsgo/cutil"
    "fmt"
    "unsafe"
)

func main() {
    for _, s := range cutil.StrArrToStringSlice(unsafe.Pointer(C.GetStrArr()), 3) {
        fmt.Println(s)
    }
}
