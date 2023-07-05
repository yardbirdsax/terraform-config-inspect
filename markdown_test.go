// Copyright (c) Josh Feierman (original copyright HashiCorp, Inc).
// SPDX-License-Identifier: MPL-2.0

package terraparse

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRenderMarkdown(t *testing.T) {
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

			fullPath := filepath.Join(path, name+".out.md")
			expected, err := ioutil.ReadFile(fullPath)
			if err != nil {
				t.Skipf("%q not found, skipping test", fullPath)
			}

			module, _ := LoadModule(path)
			if module == nil {
				t.Fatalf("result object is nil; want a real object")
			}

			var b bytes.Buffer
			buf := &b
			err = RenderMarkdown(buf, module)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(buf.String(), string(expected)); diff != "" {
				t.Errorf("actual and expected content differ:\n%s", diff)
			}
		})
	}
}
