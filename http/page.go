package http

import (
	"github.com/Cepave/agent/g"
	"github.com/toolkits/file"
	"net/http"
	"path/filepath"
	"strings"
)

func configPageRoutes() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if g.IsTrustable(r.RemoteAddr) {
			if strings.HasSuffix(r.URL.Path, "/") {
				if !file.IsExist(filepath.Join(g.Root, "/public", r.URL.Path, "index.html")) {
					http.NotFound(w, r)
					return
				}
			}
			http.FileServer(http.Dir(filepath.Join(g.Root, "/public"))).ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("no privilege"))
		}
	})

}
