package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

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
	log.Fatalln(http.ListenAndServe(":8080", h))
}
