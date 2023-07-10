// Copyright (c) Josh Feierman (original copyright HashiCorp, Inc).
// SPDX-License-Identifier: MPL-2.0

package terraparse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadModule(t *testing.T) {
	fixturesDir := "testdata"
	testDirs, err := ioutil.ReadDir(fixturesDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, info := range testDirs {
		if !info.IsDir() {
			continue
		}
		t.Run(info.Name(), func(t *testing.T) {
			name := info.Name()
			path := filepath.Join(fixturesDir, name)

			wantSrc, err := ioutil.ReadFile(filepath.Join(path, name+".out.json"))
			if err != nil {
				t.Fatalf("failed to read result file: %s", err)
			}
			var want map[string]interface{}
			err = json.Unmarshal(wantSrc, &want)
			if err != nil {
				t.Fatalf("failed to parse result file: %s", err)
			}

			gotObj, _ := LoadModule(path)
			if gotObj == nil {
				t.Fatalf("result object is nil; want a real object")
			}

			gotSrc, err := json.Marshal(gotObj)
			if err != nil {
				t.Fatalf("result is not JSON-able: %s", err)
			}
			var got map[string]interface{}
			err = json.Unmarshal(gotSrc, &got)
			if err != nil {
				t.Fatalf("failed to parse the actual result (!?): %s", err)
			}

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("wrong result\n%s", diff)
			}
		})
	}
}

func TestLoadModuleFromFilesystem(t *testing.T) {
	fixturesDir := "testdata"
	testDirs, err := ioutil.ReadDir(fixturesDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, info := range testDirs {
		if !info.IsDir() {
			continue
		}
		t.Run(info.Name(), func(t *testing.T) {
			name := info.Name()
			path := filepath.Join(fixturesDir, name)
			fs := os.DirFS(".")

			wantSrc, err := ioutil.ReadFile(filepath.Join(path, name+".out.json"))
			if err != nil {
				t.Fatalf("failed to read result file: %s", err)
			}
			var want map[string]interface{}
			err = json.Unmarshal(wantSrc, &want)
			if err != nil {
				t.Fatalf("failed to parse result file: %s", err)
			}

			gotObj, _ := LoadModuleFromFilesystem(WrapFS(fs), path)
			if gotObj == nil {
				t.Fatalf("result object is nil; want a real object")
			}

			gotSrc, err := json.Marshal(gotObj)
			if err != nil {
				t.Fatalf("result is not JSON-able: %s", err)
			}
			var got map[string]interface{}
			err = json.Unmarshal(gotSrc, &got)
			if err != nil {
				t.Fatalf("failed to parse the actual result (!?): %s", err)
			}

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("wrong result\n%s", diff)
			}
		})
	}
}

func sortModuleMapKeys[v *ModuleCall | *Attribute](in map[string]v) []string {
	out := make([]string, len(in))
	for k := range in {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func ExampleLoadModule_modulecalls() {
	got, err := LoadModule("testdata/module-calls")
	if err != nil {
		fmt.Printf("error loading file: %v", err)
	}
	if len(got.ModuleCalls) > 0 {
		fmt.Println("Module calls:\n------")
		sortedKeys := sortModuleMapKeys(got.ModuleCalls)
		for _, k := range sortedKeys {
			if mc, ok := got.ModuleCalls[k]; ok {
				fmt.Printf("%s\n", mc.Name)
				fmt.Printf("\tSource: %s\n", mc.Source)
				fmt.Println("\tAttributes:")
				sortedAttributeKeys := sortModuleMapKeys(mc.Attributes)
				for _, ak := range sortedAttributeKeys {
					if a, ok := mc.Attributes[ak]; ok {
						fmt.Printf("\t\t%s: %s\n", a.Name, a.Value.GoString())
					}
				}
				fmt.Println("---")
			}
		}
	}

	// Output:
	// Module calls:
	// ------
	// bar
	// 	Source: ./child
	// 	Attributes:
	// 		unused: cty.NumberIntVal(1)
	// ---
	// baz
	// 	Source: ../elsewhere
	// 	Attributes:
	// 		unused: cty.NumberIntVal(12)
	// ---
	// foo
	// 	Source: foo/bar/baz
	// 	Attributes:
	// 		id: cty.StringVal("data.external.something.result.id")
	// 		something: cty.StringVal("var.something")
	// 		something_else: cty.StringVal("${var.something}-2")
	// 		unused: cty.NumberIntVal(2)
	// ---
}
