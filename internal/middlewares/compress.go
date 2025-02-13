package middlewares

import (
	"compress/flate"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"slices"
	"strings"
)

func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !slices.Contains(allowContentTypes, r.Header.Get("Content-Type")) {
			next.ServeHTTP(w, r)
			return
		}

		ow := w

		if r.Header.Get("Accept-Encoding") != "" {
			cw, err := newCompressWriter(w, r.Header.Get("Accept-Encoding"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ow = cw

			defer func() {
				if err = cw.Close(); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}()
		}

		if r.Header.Get("Content-Encoding") != "" {
			cr, err := newCompressReader(r.Body, r.Header.Get("Content-Encoding"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			r.Body = cr

			defer func() {
				if err = cr.Close(); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}()
		}

		next.ServeHTTP(ow, r)
	})
}

var allowContentTypes = []string{"application/json"}

const encodingGzip = "gzip"
const encodingDeflate = "deflate"
const encodingWildcard = "*"

type compressWriter struct {
	w        http.ResponseWriter
	cw       io.WriteCloser
	encoding string
}

func (cw *compressWriter) Header() http.Header {
	return cw.w.Header()
}

func (cw *compressWriter) Write(b []byte) (int, error) {
	return cw.cw.Write(b)
}

func (cw *compressWriter) WriteHeader(statusCode int) {
	if statusCode != 300 {
		cw.w.Header().Set("Content-Encoding", cw.encoding)
	}

	cw.w.WriteHeader(statusCode)
}

func (cw *compressWriter) Close() error {
	return cw.cw.Close()
}

func newCompressWriter(w http.ResponseWriter, acceptEncoding string) (*compressWriter, error) {
	if acceptEncoding == "" {
		return nil, errors.New("Accept-Encoding is empty")
	}

	switch {
	case strings.Contains(acceptEncoding, encodingWildcard):
		fallthrough
	case strings.Contains(acceptEncoding, encodingGzip):
		cw, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			return nil, err
		}

		return &compressWriter{w: w, cw: cw, encoding: encodingGzip}, nil
	case strings.Contains(acceptEncoding, encodingDeflate):
		cw, err := flate.NewWriter(w, flate.BestCompression)
		if err != nil {
			return nil, err
		}

		return &compressWriter{w: w, cw: cw, encoding: encodingDeflate}, nil
	}

	return nil, errors.New("Accept-Encoding type is not supported")
}

type compressReader struct {
	r  io.ReadCloser
	cr io.ReadCloser
}

func (cr *compressReader) Read(b []byte) (int, error) {
	return cr.cr.Read(b)
}

func (cr *compressReader) Close() error {
	if err := cr.r.Close(); err != nil {
		return err
	}

	return cr.cr.Close()
}

func newCompressReader(r io.ReadCloser, contentEncoding string) (*compressReader, error) {
	if contentEncoding == "" {
		return nil, errors.New("Content-Encoding is empty")
	}

	switch {
	case strings.Contains(contentEncoding, encodingDeflate):
		return &compressReader{r: r, cr: flate.NewReader(r)}, nil
	case strings.Contains(contentEncoding, encodingGzip):
		cr, err := gzip.NewReader(r)
		if err != nil {
			return nil, err
		}
		return &compressReader{r: r, cr: cr}, nil
	}

	return nil, errors.New("Content-Encoding type is not supported")
}
