package terraparse

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	"github.com/zclconf/go-cty/cty/json"
)

type Attribute struct {
	*hcl.Attribute
	Value cty.Value
}

type Attributes map[string]*Attribute

func NewAttributesFromBody(body hcl.Body, file *hcl.File) (mas Attributes, diags hcl.Diagnostics) {
	attrs, attrDiags := body.JustAttributes()
	diags = append(diags, attrDiags...)
	if attrDiags.HasErrors() {
		return nil, diags
	}
	mas, attrDiags = NewAttributes(attrs, file)
	diags = append(diags, attrDiags...)
	return mas, diags
}

func NewAttributes(attrs hcl.Attributes, file *hcl.File) (mas Attributes, diags hcl.Diagnostics) {
	if len(attrs) == 0 {
		return nil, nil
	}
	reader := bytes.NewReader(file.Bytes)
	mas = make(Attributes, len(attrs))
	for k, v := range attrs {
		ma := Attribute{
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

func (ma Attributes) MarshalJSON() (data []byte, err error) {
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
	}
	out.WriteString("}")
	return out.Bytes(), err
}
