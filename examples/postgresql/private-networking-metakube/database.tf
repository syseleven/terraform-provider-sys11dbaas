resource "sys11dbaas_database_v2" "db" {
  name = "example-private-networking-vm"
  application_config = {
    instances = 1
    password  = var.database_password
    type      = "postgresql"
    version   = 17.5

    private_networking = {
      enabled            = true
      shared_subnet_cidr = "10.1.42.0/24"
      allowed_cidrs      = [local.metakube_node_subnet_cidr]
    }
  }
  service_config = {
    disksize = 5
    flavor   = "SCS-2V-4-50n"
    region   = var.region
  }
}

data "openstack_networking_router_v2" "router" {
  name = "metakube-${metakube_cluster.cluster.id}"
}

resource "openstack_networking_router_interface_v2" "routerint_db" {
  router_id = data.openstack_networking_router_v2.router.id
  subnet_id = resource.sys11dbaas_database_v2.db.application_config.private_networking.shared_subnet_id
}
