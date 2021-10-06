package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/foolin/goview"
	"github.com/gorilla/mux"
	"github.com/hoisie/mustache"
	"github.com/tomasen/realip"
)

func makeHandler(fn handleFunc, c Context) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.SetResponseWriter(w)
		c.SetRequest(r)
		status, err := fn(c)
		if status >= 400 || err != nil {
			clientIP := realip.FromRequest(r)
			log.Printf("%s: %v %s %v", r.URL.Path, status, clientIP, err)
		} else {
			clientIP := realip.FromRequest(r)
			log.Printf("%s: %v %s", r.URL.Path, status, clientIP)
		}
	})
	return handler
}

func homeHandler(c Context) (int, error) {
	c.Render("home", c.GetConfig())
	return http.StatusOK, nil
}

func renderFilename(form url.Values, c *Config) string {
	data := map[string]string{}
	for name, values := range form {
		data[name] = values[0]
	}
	filename := mustache.Render(c.Filename, data)
	return filename
}

func uploadHandler(c Context) (int, error) {
	file, header := c.Form()
	filename := renderFilename(c.GetRequest().PostForm, c.GetConfig())
	ext := path.Ext(header.Filename)
	f, err := os.Create(path.Join(c.GetStorage(), filename+ext))
	if err != nil {
		c.Render("upload", goview.M{"Message": "Fail to submit."})
		return http.StatusBadRequest, err
	}
	_, err = io.Copy(f, file)
	if err != nil {
		c.Render("upload", goview.M{"Message": "Fail to submit."})
		return http.StatusBadRequest, err
	}
	c.Render("upload", goview.M{"Message": "Submitted successfully."})
	return http.StatusOK, nil
}

func listHandler(c Context) (int, error) {
	files, err := ioutil.ReadDir(c.GetStorage())
	if err != nil {
		return http.StatusBadRequest, err
	}
	c.Render("list", goview.M{"Files": files, "Title": "All Files"})
	return http.StatusOK, nil
}

func NewHandler(c Context) http.Handler {
	r := mux.NewRouter()
	r.Handle("/", makeHandler(homeHandler, c))
	r.PathPrefix("/assets").Handler(http.FileServer(http.FS(assetsFS)))
	r.Handle("/list", makeHandler(listHandler, c))
	r.Handle("/upload", makeHandler(uploadHandler, c)).Methods("POST")
	return r
}
