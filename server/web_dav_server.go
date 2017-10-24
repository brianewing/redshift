package server

import (
	"golang.org/x/net/webdav"
	"net/http"
	"log"
	"time"
)

// todo: implement a custom webdav.FileSystem to store changes in git!
// todo: implement a custom VERSIONS method (or extend PROPFIND) to list old versions
// todo: maybe use a header on a GET to request old version?

func RunWebDavServer(addr string, directory string) error {
	webDavHandler := &webdav.Handler{
		FileSystem: webdav.Dir(directory),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			log.Println("WD", r.Method, r.RequestURI, "[" + r.RemoteAddr + "]", "|", err)
		},
	}
	server := &http.Server{Addr: addr, Handler: corsWrapper(noCacheWrapper(webDavHandler))}
	return server.ListenAndServe()
}

func slowWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "OPTIONS" { time.Sleep(4 * time.Second) }
		h.ServeHTTP(w, r)
	})
}

func corsWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Depth, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PROPFIND, MKCOL, DELETE, PUT, GET")
		h.ServeHTTP(w, r)
	})
}

func noCacheWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Pragma", "no-cache")
		h.ServeHTTP(w, r)
	})
}