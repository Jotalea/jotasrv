package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const publicPrefix = "/public/"

// newPublicHandler serves dir read-only under /public/ and hands every other
// request to next. It sits in front of the login, so the content is world-readable.
// No directory listings are produced and dotfiles are never served.
func newPublicHandler(next http.Handler, dir string) http.Handler {
	files := http.StripPrefix("/public", http.FileServer(noListingFS{http.Dir(dir)}))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public" && !strings.HasPrefix(r.URL.Path, publicPrefix) {
			next.ServeHTTP(w, r)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		for _, part := range strings.Split(r.URL.Path, "/") {
			if strings.HasPrefix(part, ".") {
				http.NotFound(w, r)
				return
			}
		}
		files.ServeHTTP(w, r)
	})
}

type noListingFS struct{ fs http.FileSystem }

func (n noListingFS) Open(name string) (http.File, error) {
	f, err := n.fs.Open(name)
	if err != nil {
		return nil, err
	}
	if st, err := f.Stat(); err == nil && st.IsDir() {
		index, err := n.fs.Open(filepath.ToSlash(filepath.Join(name, "index.html")))
		if err != nil {
			f.Close()
			return nil, os.ErrNotExist
		}
		index.Close()
	}
	return f, nil
}
