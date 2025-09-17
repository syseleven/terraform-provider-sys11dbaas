data "openstack_networking_network_v2" "ext-net" {
  name = local.external_network
}

resource "openstack_networking_network_v2" "net" {
  name           = "private-networking-vm"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet" {
  name            = "private-networking-vm"
  network_id      = openstack_networking_network_v2.net.id
  dns_nameservers = ["37.123.105.116", "37.123.105.117"]
  cidr            = local.virtual_machine_subnet_cidr
  ip_version      = 4
}

resource "openstack_networking_router_v2" "router" {
  name                = "private-networking-vm"
  admin_state_up      = true
  external_network_id = data.openstack_networking_network_v2.ext-net.id
}

resource "openstack_networking_router_interface_v2" "routerint" {
  router_id = openstack_networking_router_v2.router.id
  subnet_id = openstack_networking_subnet_v2.subnet.id
}

resource "openstack_networking_floatingip_v2" "fip" {
  pool = local.external_network
}

resource "openstack_networking_floatingip_associate_v2" "fipas" {
  floating_ip = openstack_networking_floatingip_v2.fip.address
  port_id     = openstack_networking_port_v2.port.id
}

resource "openstack_networking_port_v2" "port" {
  name               = "private-networking-vm"
  network_id         = openstack_networking_network_v2.net.id
  security_group_ids = [openstack_networking_secgroup_v2.secgroup.id]

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.subnet.id
  }
}

resource "openstack_networking_secgroup_v2" "secgroup" {
  name        = "private-networking-vm"
  description = "Allow inbound SSH/ICMP for IPv4"
}

resource "openstack_networking_secgroup_rule_v2" "sgr_ipv4_tcp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.secgroup.id
}

resource "openstack_networking_secgroup_rule_v2" "sgr_ipv4_icmp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.secgroup.id
}
