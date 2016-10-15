package fs

import "os"

func Exists(p string) bool {
    if _, err := os.Stat(p); err != nil {
        return false
    }
    return true
}


