package openapi3

import (
	"encoding/json"
	"path/filepath"
	"text/template"

	"goa.design/goa/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
	"gopkg.in/yaml.v2"
)

func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, root := range roots {
		if r, ok := root.(*expr.RootExpr); ok {
			jsonFile, yamlFile, err := openapiFiles(r)
			if err != nil {
				return files, err
			}
			if jsonFile != nil && yamlFile != nil {
				files = append(files, jsonFile, yamlFile)
			}
		}
	}
	return files, nil
}

func openapiFiles(r *expr.RootExpr) (*codegen.File, *codegen.File, error) {
	// Only create a OpenAPI specification if there are HTTP services.
	if len(r.API.HTTP.Services) == 0 {
		return nil, nil, nil
	}

	jsonPath := filepath.Join(codegen.Gendir, "http", "openapi.json")
	yamlPath := filepath.Join(codegen.Gendir, "http", "openapi.yaml")
	var (
		jsonSection *codegen.SectionTemplate
		yamlSection *codegen.SectionTemplate
	)
	{
		spec, err := NewV3(r)
		if err != nil {
			return nil, nil, err
		}
		jsonSection = &codegen.SectionTemplate{
			Name:    "openapi",
			FuncMap: template.FuncMap{"toJSON": toJSON},
			Source:  "{{ toJSON .}}",
			Data:    spec,
		}
		yamlSection = &codegen.SectionTemplate{
			Name:    "openapi",
			FuncMap: template.FuncMap{"toYAML": toYAML},
			Source:  "{{ toYAML .}}",
			Data:    spec,
		}
	}

	return &codegen.File{
			Path:             jsonPath,
			SectionTemplates: []*codegen.SectionTemplate{jsonSection},
		}, &codegen.File{
			Path:             yamlPath,
			SectionTemplates: []*codegen.SectionTemplate{yamlSection},
		}, nil

}

func toJSON(d interface{}) string {
	b, err := json.Marshal(d)
	if err != nil {
		panic("openapi: " + err.Error()) // bug
	}
	return string(b)
}

func toYAML(d interface{}) string {
	b, err := yaml.Marshal(d)
	if err != nil {
		panic("openapi: " + err.Error()) // bug
	}
	return string(b)
}
