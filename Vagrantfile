# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|

  # This base box supports the following virtualization providers:
  # - virtualbox
  # - vmware_desktop
  # - hyperv
  config.vm.box = "hashicorp/bionic64"

  config.vm.hostname = "k8ssandra-dev-vm"

  # Resource configuration for vitualbox deployment
  config.vm.provider "virtualbox" do |v|
    v.memory = 8192
    v.cpus = 4
  end

  config.vm.network "forwarded_port", guest: 9000, host: 9000

  config.vm.network "forwarded_port", guest: 1313, host: 1313

  config.vm.provision :shell, path: "vagrant/bootstrap-common.sh"

  config.vm.provision :shell, path: "vagrant/bootstrap-go.sh"

  config.vm.provision :shell, path: "vagrant/bootstrap-docker.sh"

  config.vm.provision :shell, path: "vagrant/bootstrap-kubectl.sh"

  config.vm.provision :shell, path: "vagrant/bootstrap-kind.sh"

  config.vm.provision :shell, path: "vagrant/bootstrap-helm.sh"

  config.vm.provision :shell, path: "vagrant/bootstrap-hugo.sh"

  config.vm.provision :shell, path: "vagrant/bootstrap-system-final.sh"
  
end
