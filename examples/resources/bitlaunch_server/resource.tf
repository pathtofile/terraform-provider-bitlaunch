terraform {
  required_providers {
    bitlaunch = {
      version = "0.2.0"
      source  = "pathtofile-tf/bitlaunch"
    }
  }
}

variable "token" { sensitive = true }
variable "host" { default = "DigitalOcean" }

provider "bitlaunch" {
  token = var.token
}

// Data
data "bitlaunch_image" "image" {
  host         = var.host
  distro_name  = "Ubuntu"
  version_name = "20.04 (LTS) x64"
}

data "bitlaunch_region" "region" {
  host        = var.host
  region_name = "San Francisco"
  slug        = "sfo2"
}

data "bitlaunch_size" "size" {
  host      = var.host
  cpu_count = 2
  memory_mb = 2048
}

// Resources
resource "bitlaunch_sshkey" "sshkey" {
  name    = "tf_sshkeys"
  content = var.ssh_pubkey
}

resource "bitlaunch_sever" "server" {
  host     = var.host
  name     = "tf_server"
  image_id = data.bitlaunch_image.image.id
}
