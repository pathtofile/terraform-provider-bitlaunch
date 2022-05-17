terraform {
  required_providers {
    bitlaunch = {
      source  = "pathtofile-tf/bitlaunch"
      version = "0.4.0"
    }
  }
}

variable "token" { sensitive = true }
variable "ssh_pubkey" { type = string }
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
  cpu_count = 1
  memory_mb = 1024
}

// Resources
resource "bitlaunch_sshkey" "sshkey" {
  name    = "tf_sshkeys"
  content = var.ssh_pubkey
}

resource "bitlaunch_server" "server" {
  host        = var.host
  name        = "tf_server"
  image_id    = data.bitlaunch_image.image.id
  size_id     = data.bitlaunch_size.size.id
  region_id   = data.bitlaunch_region.region.id
  ssh_keys    = [bitlaunch_sshkey.sshkey.id]
  wait_for_ip = true
}

// Outputs
output "ip_address" {
  value = bitlaunch_server.server.ipv4
}
