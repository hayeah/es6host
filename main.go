// Web server that supports ES6 modules.
//
// By convention ES6 programmers like to omit module extension names.
// All this module does is to attempt to append ".js" when looking
// up the file for a module.
package main

import (
	"flag"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
)

var rootDir string
var serveAddr string

func init() {
	flag.StringVar(&rootDir, "root", "./", "root directory for static resources")
	flag.StringVar(&serveAddr, "bind", "127.0.0.1:8888", "server address to bind to")
}

func main() {
	flag.Parse()

	err := serve()
	if err != nil {
		log.Fatalln("start server:", err)
	}
}

func serve() error {
	s := &es6Server{
		root: rootDir,
	}

	log.Println("Server listening:", serveAddr)
	log.Println("Serving content from:", rootDir)
	return http.ListenAndServe(serveAddr, s)
}

type es6Server struct {
	root string
}

func (s *es6Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.Println("request err", err)
			res.WriteHeader(500)
			res.Write([]byte(err.Error()))
		}
	}()

	filePath := path.Join(s.root, req.RequestURI)

	log.Println("GET", filePath, req.Header.Get("Content-Type"))

	exts := []string{".js"}

	i := 0

	var stat os.FileInfo
	lookupFilePath := filePath
	for {
		stat, err = os.Stat(lookupFilePath)
		if err != nil && !os.IsNotExist(err) {
			return
		}

		if os.IsNotExist(err) {
			if i >= len(exts) {
				return
			}

			lookupFilePath = filePath + exts[i]
			i++
			continue
		}

		break
	}

	// Serve directory index
	if stat.IsDir() {
		if req.RequestURI[len(req.RequestURI)-1] != '/' {
			res.Header().Set("Location", "/"+filePath+"/")
			res.WriteHeader(302)
			return
		}

		lookupFilePath = path.Join(filePath, "index.html")
	}

	if filePath != lookupFilePath {
		log.Println("\t", lookupFilePath)
	}

	f, ferr := os.Open(lookupFilePath)
	if ferr != nil {
		err = ferr
		return
	}
	defer f.Close()

	chosenExt := path.Ext(lookupFilePath)
	if chosenExt != "" {
		res.Header().Set("Content-Type", mime.TypeByExtension(chosenExt))
	}

	io.Copy(res, f)
}
