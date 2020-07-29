resource "ibm_is_vpc" "iac_iks_vpc" {
  name = "${var.project_name}-${var.environment}-vpc"
  tags = [
    "project:${var.project_name}",
    "env:${var.environment}",
    "owner:${var.owner}"
  ]
}

resource "ibm_is_public_gateway" "iac_iks_gateway" {
  name  = "${var.project_name}-${var.environment}-gateway-${format("%02s", count.index)}"
  vpc   = ibm_is_vpc.iac_iks_vpc.id
  zone  = var.vpc_zone_names[count.index]
  count = local.max_size
}

resource "ibm_is_subnet" "iac_iks_subnet" {
  count                    = local.max_size
  name                     = "${var.project_name}-${var.environment}-subnet-${format("%02s", count.index)}"
  zone                     = var.vpc_zone_names[count.index]
  vpc                      = ibm_is_vpc.iac_iks_vpc.id
  public_gateway           = ibm_is_public_gateway.iac_iks_gateway[count.index].id
  total_ipv4_address_count = 256
  resource_group           = ibm_resource_group.group.id
}

resource "ibm_is_security_group_rule" "iac_iks_security_group_rule_tcp_k8s" {
  count     = local.max_size
  group     = ibm_is_vpc.iac_iks_vpc.default_security_group
  direction = "inbound"

  tcp {
    port_min = 30000
    port_max = 32767
  }
}

// Enable to ssh to the nodes
// resource "ibm_is_security_group_rule" "iac_iks_security_group_rule_ssh_k8s" {
//   group     = ibm_is_vpc.iac_iks_vpc.default_security_group
//   direction = "inbound"
//   tcp {
//     port_min = 22
//     port_max = 22
//   }
// }

// Enable to ping to the nodes
// resource "ibm_is_security_group_rule" "iac_iks_security_group_rule_icmp_k8s" {
//   group     = ibm_is_vpc.iac_iks_vpc.default_security_group
//   direction = "inbound"
//   icmp {
//     type = 8
//   }
// }
