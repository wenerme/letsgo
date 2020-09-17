package wfs

import "os"

func Exists(p string) bool {
	if _, err := os.Stat(p); err != nil {
		return false
	}
	return true
}
func IsDir(p string) bool {
	if s, err := os.Stat(p); err != nil {
		return false
	} else {
		return s.Mode()&os.ModeDir > 0
	}
}
