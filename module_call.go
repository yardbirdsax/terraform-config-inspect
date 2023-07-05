// Copyright (c) Josh Feierman (original copyright HashiCorp, Inc).
// SPDX-License-Identifier: MPL-2.0

package terraparse

// ModuleCall represents a "module" block within a module. That is, a
// declaration of a child module from inside its parent.
type ModuleCall struct {
	Name       string     `json:"name"`
	Source     string     `json:"source"`
	Version    string     `json:"version,omitempty"`
	Attributes Attributes `json:"attributes,omitempty"`

	Pos SourcePos `json:"pos"`
}
