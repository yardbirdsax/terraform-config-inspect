variable "something" {
  type        = string
  description = "A variable."
  default     = "foo"
}

data "external" "something" {

}

module "foo" {
  source  = "foo/bar/baz"
  version = "1.0.2"

  unused         = 2
  id             = data.external.something.result.id
  something      = var.something
  something_else = "${var.something}-2"
}

module "bar" {
  source = "./child"

  unused = 1
}
