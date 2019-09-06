package openapi3_test

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	oa3 "github.com/getkin/kin-openapi/openapi3"
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/plugins/v3/openapi3"
	"goa.design/plugins/v3/openapi3/testdata"
)

var update = flag.Bool("update", false, "update golden files")

func TestOpenAPI3(t *testing.T) {
	cases := []struct {
		Name string
		DSL  func()
	}{
		{"full-dsl", testdata.FullDSL},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			root := codegen.RunDSL(t, c.DSL)
			fs, err := openapi3.Generate("", []eval.Root{root}, nil)
			if err != nil {
				t.Fatal(err)
			}
			var buf bytes.Buffer
			if err := fs[0].SectionTemplates[0].Write(&buf); err != nil {
				t.Fatal(err)
			}
			golden := filepath.Join("testdata", fmt.Sprintf("%s.json", c.Name))
			if *update {
				ioutil.WriteFile(golden, buf.Bytes(), 0644)
			}
			expected, _ := ioutil.ReadFile(golden)
			if buf.String() != string(expected) {
				t.Errorf("invalid content for %s: got\n%s\ngot vs. expected:\n%s",
					fs[0].Path, buf.String(), codegen.Diff(t, buf.String(), string(expected)))
			}

			swagger := oa3.Swagger{}
			err = swagger.UnmarshalJSON(buf.Bytes())
			if err != nil {
				t.Fatal(err)
			}

			err = swagger.Validate(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			t.Fatal("")
		})
	}
}
