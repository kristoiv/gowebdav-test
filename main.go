package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/webdav"
)

func main() {
	//fs := webdav.NewMemFS()
	cwd, _ := os.Getwd()
	fs := webdav.Dir(filepath.Join(cwd, "public"))
	ls := webdav.NewMemLS()
	l := log.New(os.Stderr, "[webdav] ", log.Lshortfile|log.LstdFlags|log.LUTC)
	h := &webdav.Handler{
		Prefix:     "/",
		FileSystem: fs,
		LockSystem: ls,
		Logger: func(r *http.Request, err error) {
			if err != nil {
				l.Printf("Error=%q from req=%#v\n", err, *r)
			} else {
				l.Printf("Request: %d %s %s\n", r.ContentLength, r.Method, r.URL)
			}
		},
	}
	log.Fatalln(http.ListenAndServe(":8080", BasicAuthMiddleware("kristoiv", "1234", h)))
}

func BasicAuthMiddleware(username, password string, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		if pair[0] != username || pair[1] != password {
			http.Error(w, "Not authorized", 401)
			return
		}

		h.ServeHTTP(w, r)
	}
}
