package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
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
	XGoogleBackend xGoogleBackend `yaml:"x-google-backend"`
}

type xGoogleBackend struct {
	Address  string
	Protocol string
}

func main() {
	for _, f := range os.Args[1:] {
		var spec swaggerSpec
		file, err := os.Open(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		err = yaml.NewDecoder(file).Decode(&spec)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		injectXGoogleBackend(&spec)

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
			s.Paths[path][op]["x-google-backend"] = xGoogleBackend{
				Address:  s.XGoogleBackend.Address,
				Protocol: s.XGoogleBackend.Protocol,
			}
		}
	}
}
