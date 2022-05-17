# BitLaunch Terraform Provider
This provider provides a basic way to create and manage [BitLaunch VMs](https://bitlaunch.io/), which can
be paid for using Bitcoin.

If you find this project useful, feel free to buy me a coffee in BTC at `16g88jxnX315CnjTDbfZ9hwuWgeSbRJdMG`

# Using
To Use, just use the `hashicorp.com/pathtofile/bitlaunch` provider.

## Example
```terraform
terraform {
  required_providers {
    bitlaunch = {
      version = "~> 0.0.1"
      source  = "hashicorp.com/pathtofile/bitlaunch"
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
```

See the full docs on Terrafor, as well as the [BitLaunch API Docs](https://developers.bitlaunch.io/reference),
for more details. Most of the API objects have a 1:1 mapping to Terraform Resources or Data sources.

# Building
## Requirements
-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.18

## Building The Provider
1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Developing the Provider
If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
