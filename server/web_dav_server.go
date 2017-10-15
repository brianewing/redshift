package server

import (
	"golang.org/x/net/webdav"
	"net/http"
	"log"
)

func RunWebDavServer(addr string, directory string) error {
	handler := &webdav.Handler{
		FileSystem: webdav.Dir(directory),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			log.Println("WD:", r.Method, r.RequestURI, "[" + r.RemoteAddr + "]", "|", err)
		},
	}
	server := &http.Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}
