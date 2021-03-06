package server

import (
	"github.com/phayes/permbits"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

// todo: implement a custom webdav.FileSystem to store changes in git!
// todo: implement a custom VERSIONS method (or extend PROPFIND) to list old versions
// todo: maybe use a header on a GET to request old version?

func RunWebDavServer(addr string, directory string, makeNewFilesExecutable bool) {
	webDavHandler := &webdav.Handler{
		FileSystem: webdav.Dir(directory),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			log.Println("DAV", r.Method, r.RequestURI, "["+r.RemoteAddr+"]", "|", err)
		},
	}

	wrappedHandler := corsWrapper(noCacheWrapper(webDavHandler))

	if makeNewFilesExecutable {
		wrappedHandler = makeNewFilesExecutableWrapper(directory, wrappedHandler)
	}

	server := &http.Server{Addr: addr, Handler: wrappedHandler}
	log.Fatalln("DAV", server.ListenAndServe())
}

func slowWrapper(wait time.Duration, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "OPTIONS" {
			time.Sleep(wait)
		}
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

func makeNewFilesExecutableWrapper(baseDir string, h http.Handler) http.Handler {
	if baseDir == "" {
		baseDir = "."
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)

		if r.Method == "PUT" {
			path := filepath.Join(baseDir, r.RequestURI)
			if perms, err := permbits.Stat(path); err == nil {
				perms.SetUserExecute(true)
				permbits.Chmod(path, perms)
			}
		}
	})
}
