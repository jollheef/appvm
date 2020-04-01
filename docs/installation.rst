Installation
============

NixOS
-----

/etc/nixos/configuration.nix::

  virtualisation.libvirtd = {
    enable = true;
    qemuVerbatimConfig = ''
      namespaces = []
      user = "${username}"
      group = "users"
    '';
  };

  users.users."${username}".extraGroups = [ ... "libvirtd" ];

shell::

  nix run -f https://code.dumpstack.io/tools/appvm/archive/master.tar.gz -c appvm

Ubuntu 19.10
------

Requirements::

  sudo apt install virt-manager curl git
  echo user = "\"$USER\"" | sudo tee -a /etc/libvirt/qemu.conf
  echo '/var/tmp/** rwlk,' | sudo tee -a /etc/apparmor.d/local/abstractions/libvirt-qemu
  curl https://nixos.org/nix/install | sh
  systemctl reboot

Use latest stable nixpkgs channel::

  nix-channel --add https://nixos.org/channels/nixos-20.03 nixpkgs
  nix-channel --update

Install appvm::

  nix-env -if https://code.dumpstack.io/tools/appvm/archive/master.tar.gz
