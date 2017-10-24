package main

/*
#include <string.h>
#include <stdlib.h>
#include <stdint.h>

int test_ints[] = {1,2,3,4};
int *get_test_ints() {
  return test_ints;
}

char *test_strs[] = {
    "Hello",
    "There",
    "中文"
};
char **get_test_strs() {
  return test_strs;
}
*/
import "C"
import (
	"fmt"
	"github.com/wenerme/letsgo/cinterop"
	"strings"
	"unsafe"
)

func main() {
	fmt.Println(strings.Join(cinterop.GoStringSlice(unsafe.Pointer(C.get_test_strs()), 3), ","))
	fmt.Printf("%v\n", cinterop.GoIntSlice(unsafe.Pointer(C.get_test_ints()), 4))
}
