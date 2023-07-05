// Copyright (c) Josh Feierman (original copyright HashiCorp, Inc).
// SPDX-License-Identifier: MPL-2.0

package terraparse

import (
	"bytes"
	// gojson "encoding/json"
	"fmt"
	"strconv"

	// "strings"

	// "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	"github.com/zclconf/go-cty/cty/json"
)

type ModuleAttribute struct {
	*hcl.Attribute
	Value cty.Value
}

type ModuleAttributes map[string]*ModuleAttribute

func NewModuleAttributes(attrs hcl.Attributes, file *hcl.File) (mas ModuleAttributes, diags hcl.Diagnostics) {
	if len(attrs) == 0 {
		return nil, nil
	}
	reader := bytes.NewReader(file.Bytes)
	mas = make(ModuleAttributes, len(attrs))
	for k, v := range attrs {
		ma := ModuleAttribute{
			Attribute: v,
		}
		var val cty.Value
		var valDiags hcl.Diagnostics
		switch expr := v.Expr.(type) {
		case *hclsyntax.LiteralValueExpr:
			val, valDiags = expr.Value(nil)
			diags = append(diags, valDiags...)
			if valDiags.HasErrors() {
				continue
			}
		default:
			rangeBytes := make([]byte, v.Expr.Range().End.Byte-v.Expr.Range().Start.Byte)
			_, err := reader.ReadAt(rangeBytes, int64(v.Expr.Range().Start.Byte))
			if err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Error reading attribute offset",
					Detail:   fmt.Sprintf("An error was encountered reading the raw offset data for attribute %q: %v", k, err),
					Subject:  &v.Range,
				})
				continue
			}
			valStr, err := strconv.Unquote(string(rangeBytes))
			if err != nil {
				valStr = string(rangeBytes)
			}
			// This is to attempt to get the value as it's proper type. It's mostly
			// useful for JSON expressions where the underyling type isn't known (it's
			// not exported from the HCL library so we can't catch it in the switch statement).
			var valToConvert interface{}
			var ct cty.Type
			func() {
				var err error
				if valToConvert, err = strconv.ParseBool(valStr); err == nil {
					ct = cty.Bool
					return
				}
				if valToConvert, err = strconv.ParseInt(valStr, 10, 64); err == nil {
					ct = cty.Number
					return
				}
				if valToConvert, err = strconv.ParseFloat(valStr, 64); err == nil {
					ct = cty.Number
					return
				}
				valToConvert = valStr
				ct = cty.String
			}()

			val, err = gocty.ToCtyValue(valToConvert, ct)
			if err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Error converting attribute value to gocty value",
					Detail:   fmt.Sprintf("An error was encountered converting the value for attribute %q: %v", k, err),
				})
			}
		}
		ma.Value = val
		mas[k] = &ma
	}
	return mas, diags
}

// ModuleCall represents a "module" block within a module. That is, a
// declaration of a child module from inside its parent.
type ModuleCall struct {
	Name       string           `json:"name"`
	Source     string           `json:"source"`
	Version    string           `json:"version,omitempty"`
	Attributes ModuleAttributes `json:"attributes,omitempty"`

	Pos SourcePos `json:"pos"`
}

func (ma ModuleAttributes) MarshalJSON() (data []byte, err error) {
	out := &bytes.Buffer{}
	out.WriteString("{")
	first := true
	for k, v := range ma {
		if !first {
			out.WriteString(",")
		}
		first = false
		out.WriteString(strconv.Quote(k) + ":")
		valData, err := json.Marshal(v.Value, v.Value.Type())
		if err != nil {
			return nil, fmt.Errorf("error marshalling field (%q): %w", k, err)
		}
		out.Write(valData)
		// switch expr := expr.(type) {
		// case *hclsyntax.ScopeTraversalExpr:
		// 	s := []string{}
		// 	for _, t := range expr.Traversal {
		// 		switch t := t.(type) {
		// 		case hcl.TraverseAttr:
		// 			s = append(s, t.Name)
		// 		case hcl.TraverseRoot:
		// 			s = append(s, t.Name)
		// 		default:
		// 			return nil, fmt.Errorf("unexpected type encountered for traversal (%q): %T", k, t)
		// 		}
		// 	}
		// 	data, err := gojson.Marshal(strings.Join(s, "."))
		// 	if err != nil {
		// 		return nil, fmt.Errorf("error marshalling field (%q): %w", k, err)
		// 	}
		// 	out.Write(data)
		// case *hclsyntax.TemplateExpr:
		// 	//s := []string{}
		// 	for _, p := range expr.Parts {
		// 		switch p := p.(type) {
		// 		default:
		// 			fmt.Printf("%#v", p)
		// 		}
		// 	}
		// default:
		// 	val, valDiags := v.Expr.Value(nil)
		// 	if valDiags.HasErrors() {
		// 		return nil, &multierror.Error{Errors: valDiags.Errs()}
		// 	}
		// 	valData, err := json.Marshal(val, val.Type())
		// 	if err != nil {
		// 		return nil, fmt.Errorf("error marshalling field (%q): %w", k, err)
		// 	}
		// 	out.Write(valData)
		// }
	}
	out.WriteString("}")
	return out.Bytes(), err
}
