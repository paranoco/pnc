package conflang

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
	"log"
	"os"
	"strings"
)

type PObjectSpec struct {
	Attrs []PAttrSpec `hcl:"attr,block"`
	BlockMaps []PBlockMapSpec `hcl:"blockmap,block"`
}

func (self *PObjectSpec) ToObjectSpec() *hcldec.ObjectSpec {
	res := hcldec.ObjectSpec{}

	for _, attr := range self.Attrs {
		res[attr.Name] = attr.ToAttrSpec()
	}

	for _, blockMap := range self.BlockMaps {
		res[blockMap.Name] = blockMap.ToBlockMapSpec()
	}

	return &res
}

type PAttrSpec struct {
	Name string `hcl:"name,label"`
	Type string `hcl:"type"`
	Required bool `hcl:"required"`
}

func stringToCtyType(s string) cty.Type {
	parts := strings.Split(s, "[")
	if len(parts) == 2 {
		typeParameter := stringToPrimitiveCtyType(strings.TrimRight(parts[1], "]"))
		return stringToCollectionType(parts[0], typeParameter)
	} else if len(parts) == 1 {
		return stringToPrimitiveCtyType(parts[0])
	} else {
		return cty.DynamicPseudoType
	}
}

func stringToCollectionType(s string, typeParameter cty.Type) cty.Type {
	switch strings.ToLower(s) {
	case "list":
	return cty.List(typeParameter)
	case "set":
		return cty.Set(typeParameter)
	case "map":
		return cty.Map(typeParameter)
	}

	return cty.DynamicPseudoType
}

func stringToPrimitiveCtyType(s string) cty.Type {
	switch strings.ToLower(s) {
	case "bool":
		return cty.Bool
	case "number":
		return cty.Number
		case "string":
			return cty.String
	}

	return cty.DynamicPseudoType
}

func (self *PAttrSpec) ToAttrSpec() *hcldec.AttrSpec {
	return &hcldec.AttrSpec{
		self.Name,
		stringToCtyType(self.Type),
		self.Required,
	}
}

type PBlockMapSpec struct {
	Name string `hcl:"name,label"`
	TypeName string `hcl:"TypeName"`
	LabelNames []string `hcl:"LabelNames"`
	Attrs []PAttrSpec `hcl:"attr,block"`
	BlockMaps []PBlockMapSpec `hcl:"blockmap,block"`
}

func (self *PBlockMapSpec) ToBlockMapSpec() *hcldec.BlockMapSpec {
	res := hcldec.BlockMapSpec{}

	res.TypeName = self.TypeName
	res.LabelNames = self.LabelNames
	nested := hcldec.ObjectSpec{}

	for _, attr := range self.Attrs {
		nested[attr.Name] = attr.ToAttrSpec()
	}

	for _, blockMap := range self.BlockMaps {
		nested[blockMap.Name] = blockMap.ToBlockMapSpec()
	}

	res.Nested = nested

	return &res
}

func GetSchemaSpecFromFile(s string) (hcldec.Spec, hcl.Diagnostics) {
	parser := hclparse.NewParser()
	schemaFile, diagnostics := parser.ParseHCLFile(s)
	if diagnostics.HasErrors() {
		return nil, diagnostics
	}
	var schema PObjectSpec
	gohcl.DecodeBody(schemaFile.Body, nil, &schema)

	spec := schema.ToObjectSpec()

	return spec, diagnostics
}

func ParseFileWithHCLSchema(hclSchemaFile string, file string) {
	parser := hclparse.NewParser()
	f, parseDiags := parser.ParseHCLFile(file)
	if parseDiags.HasErrors() {
		log.Fatal(parseDiags.Error())
	}

	ctx := NewContext()
	var variableBody ConstantBody
	decodeDiags := gohcl.DecodeBody(f.Body, ctx, &variableBody)
	if decodeDiags.HasErrors() {
		log.Fatal(decodeDiags.Error())
	}

	fmt.Printf("variableBody: %#v\n", variableBody)

	declaredVars := map[string]cty.Value{}

	for _, variable := range variableBody.Constants {
		declaredVars[variable.Name] = cty.StringVal(variable.Set)
	}

	fmt.Printf("declaredVars: %#v\n", declaredVars)

	childCtx := ctx.NewChild()
	childCtx.Variables = map[string]cty.Value{
		"var": cty.MapVal(declaredVars),
	}

	spec, diags := GetSchemaSpecFromFile(hclSchemaFile)
	if diags.HasErrors() {
		log.Fatal(diags.Error())
	}

	p := json.NewEncoder(os.Stdout)
	fmt.Printf("spec (json):\n")
	p.Encode(&spec)

	// .. parse remain using childCtx
	value, diagnostics := hcldec.Decode(variableBody.MainBody, spec, childCtx)
	if diagnostics.HasErrors() {
		log.Fatal(diagnostics.Error())
		return
	}

	a, e := (ctyjson.SimpleJSONValue{value}).MarshalJSON()
	if e != nil {
		panic(e)
	}
	fmt.Printf("json(value): %s\n", a)

	fmt.Printf("value: %#v\n", value)
}
