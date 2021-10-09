package config

import (
	"encoding/json"
	"os"
)

type Input struct {
	Name    string `json:"name"`
	Label   string `json:"label"`
	Pattern string `json:"pattern"`
}

type Form struct {
	Prefix           string  `json:"prefix"`
	Storage          string  `json:"storage"`
	Title            string  `json:"title"`
	Inputs           []Input `json:"inputs"`
	FilenameTemplate string  `json:"filenameTemplate"`
}

type Config struct {
	Title string `json:"title"`
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Forms []Form `json:"forms"`
}

func (conf *Config) GetForm(prefix string) *Form {
	for _, form := range conf.Forms {
		if form.Prefix == prefix {
			return &form
		}
	}
	return nil
}

func ParseConfig(path string) (*Config, error) {
	conf := &Config{
		Host: "0.0.0.0", // default 0.0.0.0
		Port: 8080,      // default 8080
	}
	f, err := os.Open(path)
	if err != nil {
		return conf, err
	}
	err = json.NewDecoder(f).Decode(conf)
	return conf, err
}
