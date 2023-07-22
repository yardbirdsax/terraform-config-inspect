variable "instance_type" {
  type = string
}

resource "aws_instance" "foo" {
  instance_type = var.instance_type
  cpu_options {
    core_count = 1
  }
}
