terraform {
  required_providers {
    bitlaunch = {
      version = "0.2.0"
      source  = "pathtofile-tf/bitlaunch"
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
