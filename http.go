package main

import (
	"fmt"
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
	"github.com/pkg/errors"
)

func makeHandler(fn handleFunc, c Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.SetResponseWriter(w)
		c.SetRequest(r)
		status, err := fn(c)
		username := "-"
		if r.URL.User != nil {
			if name := r.URL.User.Username(); name != "" {
				username = name
			}
		}
		uri := r.RequestURI
		if r.ProtoMajor == 2 && r.Method == "CONNECT" {
			uri = r.Host
		}
		if uri == "" {
			uri = r.URL.RequestURI()
		}
		log.Printf("%s - %s \"%s %s %s\" %d", r.Host, username, r.Method, uri, r.Proto, status)
		if err != nil {
			err = errors.Errorf("%+v", err)
			fmt.Printf("\n%+v\n\n", err)
		}
	})
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
		c.Render("500", nil)
		return http.StatusInternalServerError, err
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
