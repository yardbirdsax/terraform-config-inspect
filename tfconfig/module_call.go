// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfconfig

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty/json"
)

type ModuleAttributes hcl.Attributes

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
	for k, v := range *ma {
		out.WriteString(strconv.Quote(k) + ":")
		val, valDiags := v.Expr.Value(nil)
		if valDiags.HasErrors() {
			return nil, &multierror.Error{Errors: valDiags.Errs()}
		}
		valData, err := json.Marshal(val, val.Type())
		if err != nil {
			return nil, fmt.Errorf("error marshalling value (field: %q): %w", k, err)
		}
		out.Write(valData)
	}
	out.WriteString("}")
	data = append(data, out.Bytes()...)
	return data, err
}
