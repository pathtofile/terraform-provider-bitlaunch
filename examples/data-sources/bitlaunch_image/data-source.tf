terraform {
  required_providers {
    bitlaunch = {
      version = "~> 0.0.1"
      source  = "hashicorp.com/pathtofile/bitlaunch"
    }
  }
}

variable "token" { sensitive = true }

provider "bitlaunch" {
  token = var.token
}

data "bitlaunch_image" "example" {
  host        = "DigitalOcean"
  distro_name = "Ubuntu"
  # version_name = "20.04 (LTS) x64"
}
