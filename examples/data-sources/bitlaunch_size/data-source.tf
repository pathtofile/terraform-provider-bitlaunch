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

data "bitlaunch_size" "example" {
  host      = "DigitalOcean"
  cpu_count = 2
  memory_mb = 2048
}
