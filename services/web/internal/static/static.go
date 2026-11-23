// Package static serves the built web app: real files as they are, and
// index.html for app routes so the single-page app can handle them.
package static

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// dist holds the built web app. The checked-in copy is a placeholder; the
// web build replaces it before the image is built.
//
//go:embed all:dist
var dist embed.FS

// Dist returns the embedded build output.
func Dist() fs.FS {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err) // the embed directive guarantees dist exists
	}
	return sub
}

// Handler serves files from root. A path with no file extension that does
// not exist gets index.html (an app route). A missing path with an extension
// gets 404, so a broken asset link fails loudly instead of returning HTML.
func Handler(root fs.FS) http.Handler {
	files := http.FileServerFS(root)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if name == "" {
			name = "index.html"
		}

		if _, err := fs.Stat(root, name); err != nil {
			if path.Ext(name) != "" {
				http.NotFound(w, r)
				return
			}
			serveIndex(w, r, root)
			return
		}

		if name == "index.html" {
			serveIndex(w, r, root)
			return
		}
		if strings.HasPrefix(name, "assets/") {
			// Vite puts a content hash in every file name under assets/.
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		files.ServeHTTP(w, r)
	})
}

func serveIndex(w http.ResponseWriter, r *http.Request, root fs.FS) {
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFileFS(w, r, root, "index.html")
}
