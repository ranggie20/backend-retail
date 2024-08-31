package local

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Bucket struct {
	BaseDir string
	BaseURL string
	Prefix  string
	Handler http.Handler
}

func New(basedir, server_addr, prefix string, handler http.Handler) *Bucket {
	baseurl := "http://" + server_addr
	if strings.HasPrefix(server_addr, ":") {
		baseurl += "http://localhost" + server_addr
	}
	b := Bucket{
		BaseDir: basedir,
		BaseURL: baseurl + prefix,
		Prefix:  prefix,
		Handler: handler,
	}
	return &b
}

func (b *Bucket) Upload(filename string, file io.Reader) (string, error) {
	path := filepath.Join(b.BaseDir, filename)

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		return "", err
	}
	return url.JoinPath(b.BaseURL, filename)
}

func (b *Bucket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, b.Prefix) {
		if b.Handler != nil {
			b.Handler.ServeHTTP(w, r)
		}
		return
	}

	imgPath := r.URL.Path[len(b.Prefix):]
	path := filepath.Join(b.BaseDir, imgPath)
	contents, err := os.ReadFile(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(contents)
}
