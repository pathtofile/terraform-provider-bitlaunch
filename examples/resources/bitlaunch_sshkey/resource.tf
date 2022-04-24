terraform {
  required_providers {
    bitlaunch = {
      version = "~> 0.0.1"
      source  = "hashicorp.com/pathtofile/bitlaunch"
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
