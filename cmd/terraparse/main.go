// Copyright (c) Josh Feierman (original copyright HashiCorp, Inc).
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"encoding/json"
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/yardbirdsax/terraparse"
)

var showJSON = flag.Bool("json", false, "produce JSON-formatted output")

func main() {
	flag.Parse()

	var dir string
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	} else {
		dir = "."
	}

	module, _ := terraparse.LoadModule(dir)

	if *showJSON {
		showModuleJSON(module)
	} else {
		showModuleMarkdown(module)
	}

	if module.Diagnostics.HasErrors() {
		os.Exit(1)
	}
}

func showModuleJSON(module *terraparse.Module) {
	j, err := json.MarshalIndent(module, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error producing JSON: %s\n", err)
		os.Exit(2)
	}
	os.Stdout.Write(j)
	os.Stdout.Write([]byte{'\n'})
}

func showModuleMarkdown(module *terraparse.Module) {
	err := terraparse.RenderMarkdown(os.Stdout, module)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error rendering template: %s\n", err)
		os.Exit(2)
	}
}
