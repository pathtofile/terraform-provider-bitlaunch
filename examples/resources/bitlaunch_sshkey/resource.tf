terraform {
  required_providers {
    bitlaunch = {
      version = "0.2.0"
      source  = "pathtofile-tf/bitlaunch"
    }
  }
}

variable "token" { sensitive = true }
variable "ssh_pubkey" {}

provider "bitlaunch" {
  token = var.token
}

resource "bitlaunch_sshkey" "tf_sshkey" {
  name    = "tf_sshkeys"
  content = var.ssh_pubkey
}
