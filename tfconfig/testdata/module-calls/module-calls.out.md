
# Module `testdata/module-calls`

Provider Requirements:
* **external:** (any version)

## Input Variables
* `something` (default `"foo"`): A variable.

## Data Resources
* `data.external.something` from `external`

## Child Modules
* `bar` from `./child`
* `baz` from `../elsewhere`
* `foo` from `foo/bar/baz` (`1.0.2`)

