module "{{labelname}}" {
  # Linux
  #source = "./modules/tf-pdc-module"
  source = "git::https://stash.mynmg.com/scm/ccoe/pdc-tf.git//tf-pdc-module"
  vspherename = "{{vspherename}}"
  user = "${var.user}"
  password = "${var.pass}"
  ## Change the things below this line to what you need ##########
  clustername = "{{vcluster}}"
  vmname = "{{servername}}"
  template = "{{templatename}}"
  dcname = "{{dcname}}"
  vmNetwork = "{{networkname}}"
  datastoreclustername = "{{storagecluster}}"
  ipAddr = "{{ipaddress}}"
  gatewayIp = "{{gatewayaddress}}"
  bitMask = "{{bitmask}}"
  nmgOrg = "{{nmgorg}}"
  nmgOwner = "{{nmgowner}}"
  nmgApplication = "{{nmgapplication}}"
}

variable "pass" {}
variable "user" {}
