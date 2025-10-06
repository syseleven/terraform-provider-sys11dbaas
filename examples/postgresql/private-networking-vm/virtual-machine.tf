data "openstack_images_image_v2" "image" {
  most_recent = true
  visibility  = "public"
  properties = {
    os_distro  = "ubuntu"
    os_version = "24.04"
  }
}

resource "openstack_compute_instance_v2" "instance" {
  name            = "private-networking-vm"
  image_id        = data.openstack_images_image_v2.image.id
  flavor_name     = "SCS-1V-2-50n"
  key_pair        = openstack_compute_keypair_v2.keypair.name
  security_groups = [openstack_networking_secgroup_v2.secgroup.name]

  network {
    port = openstack_networking_port_v2.port.id
  }


  connection {
    host        = openstack_networking_floatingip_v2.fip.address
    user        = "ubuntu"
    private_key = tls_private_key.private_key.private_key_pem
  }

  provisioner "remote-exec" {
    inline = [
      "while [ ! -e /var/lib/cloud/instance/boot-finished ]; do sleep 1; echo 'Waiting for cloud-init to finish.'; done",
      "sudo apt-get update",
      "sudo apt-get install -y postgresql-client",
    ]
  }
}

# Use a Terraform managed key for provisioning
resource "tls_private_key" "private_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "local_file" "ssh_private_key" {
  content         = tls_private_key.private_key.private_key_pem
  filename        = "ssh_private.key"
  file_permission = "0400"
}

resource "openstack_compute_keypair_v2" "keypair" {
  name       = "private-networking-vm"
  public_key = "${tls_private_key.private_key.public_key_openssh}\n${var.ssh_publickey}"
}
