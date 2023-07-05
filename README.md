# terraparse

This repository contains a helper library for extracting high-level metadata
about Terraform modules from their source code. It processes only a subset
of the information Terraform itself would process; in return, it can
be broadly compatible with modules written for many different versions of
Terraform.

This project was forked from the Hashicorp project
[terraform-config-inspect](https://github.com/hashicorp/terraform-config-inspect).

```
$ go get github.com/yardbirdsax/terraparse
```

```go
import "github.com/yardbirdsax/terraparse/"

// ...

module, diags := terraparse.LoadModule(dir)

// ...
```

Due to the [Terraform v1.0 Compatibility Promises](https://www.terraform.io/docs/language/v1-compatibility-promises.html),
this library should be able to parse Terraform configurations written in
the language as defined with Terraform v1.0, although it may not immediately
expose _new_ additions to the language added during the v1.x series.

This library can also interpret valid Terraform configurations targeting
Terraform v0.10 through v0.15, although the level of detail returned may
be lower in older language versions.

## Command Line Tool

The primary way to use this repository is as a Go library, but as a convenience, it also contains a
CLI tool called `terraparse` that allows viewing module information in either a Markdown-like format
or JSON format. You can install the tool by running `go install
github.com/yardbirdsax/terraparse/cmd/terraparse`.

```sh
$ terraparse path/to/module
```
```markdown
# Module `path/to/module`

Provider Requirements:
* **null:** (any version)

## Input Variables
* `a` (default `"a default"`)
* `b` (required): The b variable

## Output Values
* `a`
* `b`: I am B

## Managed Resources
* `null_resource.a` from `null`
* `null_resource.b` from `null`
```

```sh
$ terraform-config-inspect --json path/to/module
```
```json
{
  "path": "path/to/module",
  "variables": {
    "A": {
      "name": "A",
      "default": "A default",
      "pos": {
        "filename": "path/to/module/basics.tf",
        "line": 1
      }
    },
    "B": {
      "name": "B",
      "description": "The B variable",
      "pos": {
        "filename": "path/to/module/basics.tf",
        "line": 5
      }
    }
  },
  "outputs": {
    "A": {
      "name": "A",
      "pos": {
        "filename": "path/to/module/basics.tf",
        "line": 9
      }
    },
    "B": {
      "name": "B",
      "description": "I am B",
      "pos": {
        "filename": "path/to/module/basics.tf",
        "line": 13
      }
    }
  },
  "required_providers": {
    "null": []
  },
  "managed_resources": {
    "null_resource.A": {
      "mode": "managed",
      "type": "null_resource",
      "name": "A",
      "provider": {
        "name": "null"
      },
      "pos": {
        "filename": "path/to/module/basics.tf",
        "line": 18
      }
    },
    "null_resource.B": {
      "mode": "managed",
      "type": "null_resource",
      "name": "B",
      "provider": {
        "name": "null"
      },
      "pos": {
        "filename": "path/to/module/basics.tf",
        "line": 19
      }
    }
  },
  "data_resources": {},
  "module_calls": {}
}
```

## Contributing

As with its upstream inspiration, this project is designed to allow parsing a limited set of
Terraform's own dialect. While efforts may be made to keep in sync with additions or changes to
that, there is no expectation that this will be done in a timely manner.

Bug fixes are welcome so long as they include test coverage proving they work. If you would like to
contribute an enhancement or extend the language support of this library, I would suggest opening an
issue first for discussion.
