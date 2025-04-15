# Copyright (c) HashiCorp, Inc.

variable "server_url" {
  description = "Securden server URL"
  type        = string
}

variable "authtoken" {
  description = "Authentication Token"
  type        = string
  sensitive   = true
}
