package wcompress

import (
	"compress/bzip2"
	"compress/flate"
	"compress/gzip"
	"context"
	"io"
)

func Decompress(reader io.Reader, file string) (fn string, r io.Reader, err error) {
	return DecompressWithContext(reader, file, context.Background())
}
func DecompressWithContext(reader io.Reader, file string, ctx context.Context) (fn string, r io.Reader, err error) {
	fn = file
	r = reader
	if !IsCompressed(fn) {
		return
	}
	return
}
func DecompressRecursive(reader io.Reader, file string) (fn string, r io.Reader, err error) {
	return DecompressRecursiveWithContext(reader, file, context.Background())
}
func DecompressRecursiveWithContext(reader io.Reader, file string, ctx context.Context) (fn string, r io.Reader, err error) {
	fn = file
	r = reader
	for IsCompressed(fn) {
		fn, r, err = DecompressWithContext(r, fn, ctx)
		if err != nil {
			return
		}
	}
	return
}

func init() {
	// .bz2, .tar.bz2, .tbz2, .tb2
	// application/x-bzip
	//
	// .zip
	// application/zip

	Registry(Compressor{
		Name: "gzip",
		Ext:  []string{".gz"},
		Decompress: func(ctx context.Context, reader io.Reader) (io.Reader, error) {
			return gzip.NewReader(reader)
		},
		Compress: func(ctx context.Context, writer io.Writer) (w io.Writer, err error) {
			w = gzip.NewWriter(writer)
			return
		},
	})

	Registry(Compressor{
		Name: "bz2",
		Ext:  []string{".bz2"},
		Decompress: func(ctx context.Context, reader io.Reader) (r io.Reader, err error) {
			r = bzip2.NewReader(reader)
			return
		},
	})
	Registry(Compressor{
		Name: "deflate",
		Ext:  []string{".z"},
		Decompress: func(ctx context.Context, reader io.Reader) (r io.Reader, err error) {
			r = flate.NewReader(reader)
			return
		},
		Compress: func(ctx context.Context, writer io.Writer) (w io.Writer, err error) {
			w, err = flate.NewWriter(writer, flate.DefaultCompression)
			return
		},
	})
}
