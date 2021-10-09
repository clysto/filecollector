package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Masterminds/sprig"
	"github.com/clysto/filecollector/config"
	fchttp "github.com/clysto/filecollector/http"
	"github.com/clysto/filecollector/version"
	"github.com/foolin/goview"
)

//go:embed assets
var assetsFS embed.FS

//go:embed templates
var templatesFS embed.FS

func embeddedFH(config goview.Config, tmpl string) (string, error) {
	path := filepath.Join(config.Root, tmpl)
	bytes, err := templatesFS.ReadFile(path + config.Extension)
	return string(bytes), err
}

func funcMap() template.FuncMap {
	funcs := sprig.FuncMap()
	funcs["version"] = func() string {
		return version.Version
	}
	return funcs
}

func main() {
	configPath := flag.String("c", "filecollector.json", "config file path")
	showVersion := flag.Bool("v", false, "show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("File Browser v%s/%s\n", version.Version, version.CommitSHA)
		return
	}

	gv := goview.New(goview.Config{
		Root:      "templates",
		Extension: ".gohtml",
		Master:    "base",
		Funcs:     funcMap(),
	})

	gv.SetFileHandler(embeddedFH)

	conf, err := config.ParseConfig(*configPath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	r, err := fchttp.NewHandler(conf, gv)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	http.Handle("/assets/", http.StripPrefix("", http.FileServer(http.FS(assetsFS))))
	http.Handle("/", r)

	log.Printf("listening on http://%s:%d\n", conf.Host, conf.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Host, conf.Port), nil)
}
