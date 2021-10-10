package http

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/clysto/filecollector/config"
	"github.com/foolin/goview"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hoisie/mustache"
)

type Handler struct {
	conf       *config.Config
	router     chi.Router
	viewEngine *goview.ViewEngine
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) formPage(w http.ResponseWriter, r *http.Request) {
	prefix := chi.URLParam(r, "prefix")
	form := h.conf.GetForm(prefix)
	if form == nil {
		h.notFoundPage(w, r)
		return
	}
	h.Render(w, "form", form)
}

func (h *Handler) filesPage(w http.ResponseWriter, r *http.Request) {
	prefix := chi.URLParam(r, "prefix")
	form := h.conf.GetForm(prefix)
	if form == nil {
		h.notFoundPage(w, r)
		return
	}
	files, err := ioutil.ReadDir(form.Storage)
	if err != nil {
		h.RenderWithStatusCode(w, http.StatusInternalServerError, "500", nil)
		return
	}
	h.Render(w, "files", goview.M{"Files": files, "Title": "All Files", "Form": form})
}

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	h.Render(w, "home", goview.M{"Title": h.conf.Title, "Forms": h.conf.Forms})
}

func (h *Handler) uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header := MultipartForm(r)
	prefix := chi.URLParam(r, "prefix")
	form := h.conf.GetForm(prefix)

	if form == nil {
		h.notFoundPage(w, r)
		return
	}

	if form.Deadline != nil && form.Deadline.Before(time.Now()) {
		// 截止上传
		h.RenderWithStatusCode(w, http.StatusBadRequest, "upload", goview.M{"Message": "File submission is closed."})
		return
	}

	data := map[string]string{}
	for name, values := range r.PostForm {
		data[name] = values[0]
	}
	filename := mustache.Render(form.FilenameTemplate, data)
	ext := path.Ext(header.Filename)
	f, err := os.Create(path.Join(form.Storage, filename+ext))
	if err != nil {
		h.RenderWithStatusCode(w, http.StatusBadRequest, "upload", goview.M{"Message": "Fail to submit."})
		return
	}
	_, err = io.Copy(f, file)
	if err != nil {
		h.RenderWithStatusCode(w, http.StatusBadRequest, "upload", goview.M{"Message": "Fail to submit."})
		return
	}
	h.Render(w, "upload", goview.M{"Message": "Submitted successfully."})
}

func (h *Handler) notFoundPage(w http.ResponseWriter, r *http.Request) {
	h.RenderWithStatusCode(w, http.StatusNotFound, "404", nil)
}

func (h *Handler) Render(w http.ResponseWriter, name string, data interface{}) {
	h.RenderWithStatusCode(w, http.StatusOK, name, data)
}

func (h *Handler) RenderWithStatusCode(w http.ResponseWriter, statusCode int, name string, data interface{}) {
	err := h.viewEngine.Render(w, statusCode, name, data)
	if err != nil {
		panic(err)
	}
}

func MultipartForm(r *http.Request) (multipart.File, *multipart.FileHeader) {
	_ = r.ParseMultipartForm(10 << 20)
	file, header, _ := r.FormFile("file")
	return file, header
}

func NewHandler(conf *config.Config, viewEngine *goview.ViewEngine) (http.Handler, error) {
	r := chi.NewRouter()
	h := &Handler{
		conf:       conf,
		router:     r,
		viewEngine: viewEngine,
	}
	r.Use(middleware.Logger)
	r.NotFound(h.notFoundPage)
	r.Get("/", h.homePage)
	r.Get("/{prefix}", h.formPage)
	r.Post("/{prefix}/upload", h.uploadHandler)
	r.Get("/{prefix}/files", h.filesPage)
	return h, nil
}
