package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

type swaggerSpec struct {
	Swagger string
	Info    struct {
		Title   string
		Version string
	}
	Host           string
	Schemes        []string
	Produces       []string
	Paths          map[string]map[string]map[string]interface{}
	XGoogleBackend *xGoogleBackend `yaml:"x-google-backend,omitempty"`
}

type xGoogleBackend struct {
	Address         string
	Protocol        string
	PathTranslation string `yaml:"path_translation"`
}

func main() {
	var spec swaggerSpec
	var files []io.Reader

	if len(os.Args) > 1 {
		for _, f := range os.Args[1:] {
			file, err := os.Open(f)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			files = append(files, file)
		}
	} else {
		files = []io.Reader{os.Stdin}
	}

	for _, f := range files {
		err := yaml.NewDecoder(f).Decode(&spec)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		injectXGoogleBackend(&spec)
		spec.XGoogleBackend = nil

		err = yaml.NewEncoder(os.Stdout).Encode(spec)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func injectXGoogleBackend(s *swaggerSpec) {
	for path, ops := range s.Paths {
		for op, _ := range ops {
			if len(s.Paths[path][op]) > 0 {
				s.Paths[path][op]["x-google-backend"] = &xGoogleBackend{
					Address:         s.XGoogleBackend.Address,
					Protocol:        s.XGoogleBackend.Protocol,
					PathTranslation: s.XGoogleBackend.PathTranslation,
				}
			}
		}
	}
}
