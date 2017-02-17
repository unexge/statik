package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

// Handler struct is request handler for statik.
type Handler struct {
	config *Config
}

// NewHandler func creates a new handler.
func NewHandler(config *Config) (*Handler, error) {
	return &Handler{config}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := normalizePath(r.URL.String())

	if code, err := validatePath(p); err != nil {
		http.Error(w, err.Error(), code)

		return
	}

	config := h.config.GetConfigForPath(p)

	if err := handleServerPushs(w, config); err != nil {
		http.Error(w, "failed to push", http.StatusInternalServerError)

		return
	}

	fullPath := h.resolveFullPath(p)

	// TODO: serve from cache.

	f, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
		} else {
			http.Error(w, "server error", http.StatusInternalServerError)
		}

		return
	}
	defer f.Close()

	s, err := f.Stat()
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)

		return
	}

	if s.IsDir() {
		http.Error(w, "forbidden", http.StatusForbidden)

		return
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)

		return
	}

	buf := bytes.NewBuffer(b)

	if config.Gzip && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		if err = h.compressWithGzip(w, buf); err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)

			return
		}
	}

	fmt.Println(string(buf.Bytes()))

	http.ServeContent(w, r, s.Name(), s.ModTime(), bytes.NewReader(buf.Bytes()))
}

func (h *Handler) resolveFullPath(p string) string {
	return path.Join(h.config.Root, p)
}

func (h *Handler) compressWithGzip(w http.ResponseWriter, buf *bytes.Buffer) error {
	b := buf.Bytes()
	buf.Reset()

	writer := gzip.NewWriter(buf)

	_, err := writer.Write(b)
	if err != nil {
		return err
	}

	if err = writer.Close(); err != nil {
		return err
	}

	w.Header().Set("Content-Encoding", "gzip")

	return nil
}

func normalizePath(path string) string {
	if strings.HasSuffix(path, "/") {
		path += "index.html"
	}

	return path
}

func validatePath(path string) (code int, err error) {
	if strings.Contains(path, "..") {
		err = fmt.Errorf("invalid URL path")
		code = http.StatusBadRequest
	}

	return
}

func handleServerPushs(w http.ResponseWriter, c *FileConfig) error {
	if len(c.Push) == 0 {
		return nil
	}

	pusher, ok := w.(http.Pusher)
	if !ok {
		return nil
	}

	for _, p := range c.Push {
		if err := pusher.Push(p, nil); err != nil {
			return fmt.Errorf("failed to push %q: %q", p, err)
		}
	}

	return nil
}
