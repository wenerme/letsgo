package wfs

import "path/filepath"

func Ext(path string) (ext string, fn string) {
	ext = filepath.Ext(path)
	fn = path[:len(path)-len(ext)]
	return
}
