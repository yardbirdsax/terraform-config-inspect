
# Module `testdata/resource-with-inputs`

Provider Requirements:
* **aws:** (any version)
* **notaws:** (any version)

## Input Variables
* `instance_type` (required)

## Managed Resources
* `aws_instance.foo` from `aws`
* `aws_instance.json_bar` from `notaws`
* `aws_instance.json_baz` from `aws`

