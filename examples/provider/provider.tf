# Copyright (c) HashiCorp, Inc.

terraform {
  required_version = ">= 1.8"
  required_providers {
    securden = {
      source  = "terraform.local/local/securden"
      version = "0.0.1"
    }
  }
}

provider "securden" {
  authtoken      = var.authtoken
  server_url     = var.server_url
  server_timeout = 20
}
