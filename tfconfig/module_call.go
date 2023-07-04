// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfconfig

import (
	"bytes"
	gojson "encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty/json"
)

type ModuleAttributes struct {
	hcl.Attributes
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

func (ma *ModuleAttributes) MarshalJSON() (data []byte, err error) {
	out := &bytes.Buffer{}
	out.WriteString("{")
	first := true
	for k, v := range ma.Attributes {
		if !first {
			out.WriteString(",")
		}
		first = false
		out.WriteString(strconv.Quote(k) + ":")
		expr := v.Expr
		switch expr := expr.(type) {
		case *hclsyntax.ScopeTraversalExpr:
			s := []string{}
			for _, t := range expr.Traversal {
				switch t := t.(type) {
				case hcl.TraverseAttr:
					s = append(s, t.Name)
				case hcl.TraverseRoot:
					s = append(s, t.Name)
				default:
					return nil, fmt.Errorf("unexpected type encountered for traversal (%q): %T", k, t)
				}
			}
			data, err := gojson.Marshal(strings.Join(s, "."))
			if err != nil {
				return nil, fmt.Errorf("error marshalling field (%q): %w", k, err)
			}
			out.Write(data)
		default:
			val, valDiags := v.Expr.Value(nil)
			if valDiags.HasErrors() {
				return nil, &multierror.Error{Errors: valDiags.Errs()}
			}
			valData, err := json.Marshal(val, val.Type())
			if err != nil {
				return nil, fmt.Errorf("error marshalling field (%q): %w", k, err)
			}
			out.Write(valData)
		}
	}
	out.WriteString("}")
	return out.Bytes(), err
}
