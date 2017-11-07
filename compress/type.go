package wcompress

import (
	"context"
	"io"
	"path/filepath"
)

type Compressor struct {
	Name       string
	Ext        []string
	Decompress func(ctx context.Context, reader io.Reader) (r io.Reader, err error)
	Compress   func(ctx context.Context, writer io.Writer) (w io.Writer, err error)
}

var comopressors []*Compressor

func Registry(compressor Compressor) {
	comopressors = append(comopressors, &compressor)
}
func FindCompressorByExt(ext string) *Compressor {
	for _, v := range comopressors {
		for _, e := range v.Ext {
			if ext == e {
				return v
			}
		}
	}
	return nil
}
func IsCompressed(fn string) bool {
	return FindCompressorByExt(filepath.Ext(fn)) != nil
}
