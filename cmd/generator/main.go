package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/vektah/gqlparser/v2/ast"
	base_generated "gitlab.slade360emr.com/go/base/graph/generated"
	uploads_generated "gitlab.slade360emr.com/go/uploads/graph/generated"
)

func generate() error {
	var cfg *config.Config
	var err error

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("can't get current working dir")
	}

	path := fmt.Sprintf("%s/gqlgen.yml", filepath.Dir(cwd))
	_, err = ioutil.ReadFile(path) /* #nosec */
	if err != nil {
		log.Fatalf("can't find config file at %s: %s", path, err)
	}

	cfg, err = config.LoadConfig(path)
	if err != nil {
		log.Fatalf("can't load config from %s: %s", path, err)
	}

	sources := []*ast.Source{}
	sources = append(sources, base_generated.Sources()...)
	sources = append(sources, uploads_generated.Sources()...)

	baseSeen := false    // "graph/base.graphql" seen
	uploadsSeen := false // "uploads.graphql" seen
	for _, src := range sources {
		// append base.graphql once
		if src.Name == "graph/base.graphql" {
			if !baseSeen {
				cfg.Sources = append(cfg.Sources, src)
				baseSeen = true
				fmt.Println("added graph/base.graphql and marked it as seen")
			}
			continue
		}

		// append uploads.graphql once
		if src.Name == "uploads.graphql" {
			if !uploadsSeen {
				cfg.Sources = append(cfg.Sources, src)
				uploadsSeen = true
				fmt.Println("added uploads.graphql and marked it as seen")
			}
			continue
		}

		// append all other sources apart from federation directives
		if src.Name != "federation/directives.graphql" {
			cfg.Sources = append(cfg.Sources, src)
		}
	}

	if err = api.Generate(cfg); err != nil {
		return err
	}
	return nil
}

func main() {
	err := generate()
	if err != nil {
		log.Printf("failed to generate: %s", err)
		os.Exit(3)
	}
}
