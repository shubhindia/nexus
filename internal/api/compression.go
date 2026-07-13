package api

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	"github.com/klauspost/compress/zstd"
)

func RequestBody(r *http.Request) (io.ReadCloser, error) {
	switch r.Header.Get("Content-Encoding") {

	case "", "identity":
		return r.Body, nil

	case "gzip":
		gr, err := gzip.NewReader(r.Body)
		if err != nil {
			return nil, err
		}
		return gr, nil

	case "zstd":
		zr, err := zstd.NewReader(r.Body)
		if err != nil {
			return nil, err
		}

		return io.NopCloser(zr.IOReadCloser()), nil

	default:
		return nil, fmt.Errorf(
			"unsupported content encoding: %s",
			r.Header.Get("Content-Encoding"),
		)
	}
}
