package conflang

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/nsf/jsondiff"
	"github.com/zclconf/go-cty/cty"
	"testing"
)

func TestNewSpecToJson(t *testing.T) {
	spec, diagnostic := GetSchemaSpecFromFile("schema.hcl")
	if diagnostic.HasErrors() {
		t.Fatal(diagnostic.Error())
	}

	spec_old := hcldec.ObjectSpec{
		"io_mode": &hcldec.AttrSpec{
			Name: "io_mode",
			Type: cty.String,
		},
		"services": &hcldec.BlockMapSpec{
			TypeName:   "service",
			LabelNames: []string{"type", "label"},
			Nested: hcldec.ObjectSpec{
				"listen_addr": &hcldec.AttrSpec{
					Name:     "listen_addr",
					Type:     cty.String,
					Required: true,
				},
				"processes": &hcldec.BlockMapSpec{
					TypeName:   "process",
					LabelNames: []string{"name"},
					Nested: hcldec.ObjectSpec{
						"command": &hcldec.AttrSpec{
							Name:     "command",
							Type:     cty.List(cty.String),
							Required: true,
						},
					},
				},
			},
		},
	}

	s, err := json.MarshalIndent(&spec, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	s2, err := json.MarshalIndent(&spec_old, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	diffOpts := jsondiff.DefaultConsoleOptions()
	res, diff := jsondiff.Compare(s2, s, &diffOpts)

	if res != jsondiff.FullMatch {
		t.Errorf("the specs aren't identical: %s", diff)
	}
}

func TestDbSchema(t *testing.T) {
	v, err := ParseConfig("dbschema.hcl", "dbschema_sample.hcl")
	if err.HasErrors() {
		t.Fatal(err.Error())
	}

	fmt.Printf("%#v\n", v)
}
