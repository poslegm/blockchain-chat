package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"mime"
	"net/http"
	"path"
	"path/filepath"
)

func createPageHandler(rootDir string) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		file := mux.Vars(req)["file"]
		if file == "" {
			file = "index.html"
		}
		resp.Header().Add("Content-Type", "text/html")

		content, err := ioutil.ReadFile(path.Join(rootDir, file))
		if err != nil || content == nil {
			fmt.Println(err.Error())
			resp.WriteHeader(404)
			resp.Write([]byte{})
			return
		}

		resp.Write(content)
	}
}

func createPathHandler(rootDir string) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		dir := mux.Vars(req)["dir"]
		file := mux.Vars(req)["path"]
		ext := filepath.Ext(file)
		resp.Header().Add("Content-Type", mime.TypeByExtension(ext))

		content, err := ioutil.ReadFile(path.Join(rootDir, dir, file))
		if err != nil || content == nil {
			fmt.Println(err.Error())
			resp.WriteHeader(404)
			resp.Write([]byte{})
			return
		}

		resp.Write(content)
	}
}
