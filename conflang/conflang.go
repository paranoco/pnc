package conflang

import (
	"fmt"
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

type Constant struct {
	Name    string `hcl:"name,label"`
	Set string `hcl:"set"`
}

//type Locals struct {
//	Name string `hcl:"set"`
//}

type VStep1 struct {
	Name string `hcl:"name,label"`
	Type string `hcl:"type"`
	Remain  hcl.Body   `hcl:",remain"`
}

type VStep2 struct {
	//Name string `hcl:"name,label"`
	//Type string `hcl:"type"`
	Default string `hcl:"default"`
}

type VariableStep1 struct {
	Variable []VStep1 `hcl:"variable,block"`
	//Type string `hcl:"type"`
	//Default string `hcl:"default"`
	Remain  hcl.Body   `hcl:",remain"`
}

type ConstantBody struct {
	Constants []Constant `hcl:"constant,block"`
	MainBody  hcl.Body   `hcl:",remain"`
	//Locals Locals `hcl:"locals,block"`
}

type fooInstance struct {
	A string   `hcl:"a"`
	B []string `hcl:"b"`
}

func ParseConfig(schemaFile string, dataFile string) (cty.Value, hcl.Diagnostics) {
	spec, err := GetSchemaSpecFromFile(schemaFile)
	if err.HasErrors() {
		return cty.UnknownVal(cty.DynamicPseudoType), err
	}

	parser := hclparse.NewParser()
	f, err := parser.ParseHCLFile(dataFile)
	if err.HasErrors() {
		return cty.UnknownVal(hcldec.ImpliedType(spec)), err
	}

	ctx := NewContext()
	ctx.Variables["path"] = cty.ObjectVal(map[string]cty.Value{
		"root":    cty.StringVal("TODO"),
		"module":  cty.StringVal("TODO"),
		"current": cty.StringVal("TODO"),
	})

	var constantBody1 VariableStep1
	err = gohcl.DecodeBody(f.Body, ctx, &constantBody1)
	fmt.Printf("var: %#v\n", constantBody1)
	fmt.Printf("err: %#v\n", err.Error())
	//var cb2 VStep2
	//err = gohcl.DecodeBody(constantBody1.Variable[0].Remain, ctx, &cb2)
	//fmt.Printf("body: %#v\n", constantBody1.Variable[0].Remain)
	//fmt.Printf("var: %#v\n", cb2)
	//fmt.Printf("err: %#v\n", err.Error())

	var constantBody ConstantBody
	for _, traversal := range hcldec.Variables(constantBody1.Remain, spec) {
		fmt.Printf("traversal: (%s)  %#v\n", traversal.RootName(), traversal)
	}


	err = gohcl.DecodeBody(constantBody1.Remain, ctx, &constantBody)
	if err.HasErrors() {
		return cty.UnknownVal(hcldec.ImpliedType(spec)), err
	}
	// ... parse remain using a context enriched with variables
	childCtx := CreateContextWithConstants(constantBody, ctx)
	value, err := hcldec.Decode(constantBody.MainBody, spec, childCtx)
	if err.HasErrors() {
		return cty.UnknownVal(hcldec.ImpliedType(spec)), err
	}

	return value, nil
}

func CreateContextWithConstants(constantBody ConstantBody, ctx *hcl.EvalContext) *hcl.EvalContext {
	declaredVars := map[string]cty.Value{}
	for _, x := range constantBody.Constants {
		declaredVars[x.Name] = cty.StringVal(x.Set)
	}

	childCtx := ctx.NewChild()
	var varValue cty.Value
	if len(declaredVars) > 0 {
		varValue = cty.MapVal(declaredVars)
	} else {
		varValue = cty.MapValEmpty(cty.String)
	}

	childCtx.Variables = map[string]cty.Value{
		"const": varValue,
	}
	return childCtx
}
