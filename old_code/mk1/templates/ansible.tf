module "ansible" {
  # Linux
  source = "./modules/tf-pdc-module"
  vspherename = "PLPVCCORP01"
  user = "${var.user}"
  password = "${var.pass}"
  ## Change the things below this line to what you need ##########
  clustername = "RHEL Corp"
  vmname = "plplansp01"
  template = "rhel7template"
  dcname = "PDC Corp"
  #dcname = "PDC SDE"
  vmNetwork = "InfraTools"
  #datastorename = "VMAX-Corp"
  datastoreclustername = "VMAX-Corp"
  ipAddr = "172.18.35.48"
  gatewayIp = "172.18.35.1"
  bitMask = "24"
  nmgOrg = "Infrastructure"
  nmgOwner = "unixteam@neimanmarcus.com"
  nmgApplication = "reposync"
}
