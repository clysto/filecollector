package main

import (
	"log"
	"mime/multipart"
	"net/http"

	"github.com/foolin/goview"
)

type Context interface {
	Render(filename string, data interface{})
	SetResponseWriter(w http.ResponseWriter)
	SetRequest(r *http.Request)
	GetRequest() *http.Request
	GetConfig() *Config
	Form() (multipart.File, *multipart.FileHeader)
	GetStorage() string
}

type context struct {
	responseWriter http.ResponseWriter
	request        *http.Request
	template       *goview.ViewEngine
	config         *Config
}

func (c *context) GetRequest() *http.Request {
	return c.request
}

func (c *context) Form() (multipart.File, *multipart.FileHeader) {
	_ = c.request.ParseMultipartForm(10 << 20)
	file, header, _ := c.request.FormFile("file")
	return file, header
}

func (c *context) Render(name string, data interface{}) {
	err := c.template.Render(c.responseWriter, http.StatusOK, name, data)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *context) SetResponseWriter(w http.ResponseWriter) {
	c.responseWriter = w
}

func (c *context) SetRequest(r *http.Request) {
	c.request = r
}

func (c *context) GetConfig() *Config {
	return c.config
}

func (c *context) GetStorage() string {
	return c.config.Storage
}
