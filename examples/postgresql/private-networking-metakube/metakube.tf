data "metakube_k8s_version" "latest" {}

data "environment_variables" "project_id" {
  filter    = "SYS11DBAAS_PROJECT"
  sensitive = true
}

resource "metakube_cluster" "cluster" {
  name       = "dbaas-private-networking"
  dc_name    = "syseleven-${var.region}"
  project_id = base64decode(data.environment_variables.project_id.items["SYS11DBAAS_PROJECT"])

  spec {
    version = data.metakube_k8s_version.latest.version

    cloud {
      openstack {
        application_credentials {}
        subnet_cidr = local.metakube_node_subnet_cidr
      }
    }
  }
}

resource "metakube_node_deployment" "node_deployment" {
  cluster_id = metakube_cluster.cluster.id
  project_id = base64decode(data.environment_variables.project_id.items["SYS11DBAAS_PROJECT"])

  spec {
    replicas = 1
    template {
      cloud {
        openstack {
          image  = "Flatcar Stable"
          flavor = "SCS-4V-8-50n"
        }
      }
      operating_system {
        flatcar {
          disable_auto_update = true
        }
      }
    }
  }
}

resource "local_file" "kubeconfig" {
  content         = metakube_cluster.cluster.kube_config
  file_permission = "0600"
  filename        = "${path.module}/kubeconfig"
}
