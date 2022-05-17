terraform {
  required_providers {
    bitlaunch = {
      version = "0.4.0"
      source  = "pathtofile-tf/bitlaunch"
    }
  }
}

variable "token" { sensitive = true }

provider "bitlaunch" {
  token = var.token
}

data "bitlaunch_region" "example" {
  host        = "DigitalOcean"
  region_name = "New York"
  slug        = "nyc1"
}
