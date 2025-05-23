Installation
============

NixOS
-----

  virtualisation.appvm = {
    enable = true;
    user = "${username}";
  };

Ubuntu 20.04
------

Requirements::

  sudo apt install virt-manager curl git
  echo user = "\"$USER\"" | sudo tee -a /etc/libvirt/qemu.conf
  echo '/var/tmp/** rwlk,' | sudo tee -a /etc/apparmor.d/local/abstractions/libvirt-qemu
  curl -L https://nixos.org/nix/install | sh
  systemctl reboot

Use latest stable nixpkgs channel::

  nix-channel --add https://nixos.org/channels/nixos-20.03 nixpkgs
  nix-channel --update

Install appvm::

  nix-env -if https://code.dumpstack.io/tools/appvm/archive/master.tar.gz
