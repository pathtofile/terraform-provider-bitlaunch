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

data "bitlaunch_size" "example" {
  host      = "DigitalOcean"
  cpu_count = 2
  memory_mb = 2048
}
