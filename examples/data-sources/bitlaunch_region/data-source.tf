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

data "bitlaunch_region" "example" {
  host        = "DigitalOcean"
  region_name = "New York"
  slug        = "nyc1"
}
