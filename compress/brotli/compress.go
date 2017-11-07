package wbrotli

import (
	"context"
	"github.com/wenerme/letsgo/compress"
	"gopkg.in/kothar/brotli-go.v0/dec"
	"gopkg.in/kothar/brotli-go.v0/enc"
	"io"
)

func init() {
	wcompress.Registry(wcompress.Compressor{
		Name: "brotli",
		Ext:  []string{".br"},
		Decompress: func(ctx context.Context, reader io.Reader) (r io.Reader, err error) {
			r = dec.NewBrotliReader(reader)
			return
		},
		Compress: func(ctx context.Context, writer io.Writer) (w io.Writer, err error) {
			params := enc.NewBrotliParams()
			w = enc.NewBrotliWriter(params, writer)
			return
		},
	})
}
