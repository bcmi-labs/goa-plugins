package openapi3_test

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

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
		{"simple-dsl", testdata.SimpleDSL},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			root := codegen.RunDSL(t, c.DSL)
			fs, err := openapi3.Generate("", []eval.Root{root}, nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(fs) == 0 {
				t.Fatalf("got 0 files, expected 1")
			}
			if len(fs[0].SectionTemplates) == 0 {
				t.Fatalf("got 0 sections, expected 1")
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
		})
	}
}
