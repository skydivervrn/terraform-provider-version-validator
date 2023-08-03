terraform {
  required_providers {
    version-validator = {
      source  = "skydivervrn/version-validator"
      version = "0.0.5"
    }
  }
}

provider "version-validator" {}

locals {
  required_version = ">3.3.41"
}

variable "current_version" {
  default = "3.3.43"
}

data "version_validator" "example" {
  provider         = version-validator
  current_version  = var.current_version
  required_version = local.required_version
}