package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Masterminds/sprig"

	"github.com/foolin/goview"
)

//go:embed assets
var assetsFS embed.FS

//go:embed templates
var templatesFS embed.FS

type Input struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

type Config struct {
	Host     string  `json:"host"`
	Port     int     `json:"port"`
	Storage  string  `json:"storage"`
	Title    string  `json:"title"`
	Inputs   []Input `json:"inputs"`
	Filename string  `json:"filename"`
}

type handleFunc func(c Context) (int, error)

func embeddedFH(config goview.Config, tmpl string) (string, error) {
	path := filepath.Join(config.Root, tmpl)
	bytes, err := templatesFS.ReadFile(path + config.Extension)
	return string(bytes), err
}

func main() {
	configPath := flag.String("c", "filecollector.json", "config file path")
	flag.Parse()

	// 初始化模板引擎
	gv := goview.New(goview.Config{
		Root:      "templates",
		Extension: ".gohtml",
		Master:    "base",
		Funcs:     sprig.FuncMap(),
	})

	gv.SetFileHandler(embeddedFH)

	// 加载配置
	f, err := os.Open(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	config := &Config{
		Host:    "0.0.0.0",
		Port:    8080,
		Storage: "files",
	}
	err = json.NewDecoder(f).Decode(config)
	if err != nil {
		log.Fatal(err)
	}

	// 上下文
	c := &context{
		template: gv,
		config:   config,
	}

	r := NewHandler(c)
	http.Handle("/", r)

	log.Printf("listening on http://%s:%d\n", config.Host, config.Port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
