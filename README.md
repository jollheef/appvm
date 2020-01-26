[![Documentation Status](https://readthedocs.org/projects/appvm/badge/?version=latest)](https://appvm.readthedocs.io/en/latest/?badge=latest)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=R8W2UQPZ5X5JE&source=url)
[![Donate](https://img.shields.io/badge/Donate-BitCoin-green.svg)](https://blockchair.com/bitcoin/address/bc1q23fyuq7kmngrgqgp6yq9hk8a5q460f39m8nv87)

# Nix application VMs: security through virtualization

Simple application VMs (hypervisor-based sandbox) based on Nix package manager.

Uses one **read-only** /nix directory for all appvms. So creating a new appvm (but not first) is just about one minute.

![appvm screenshot](screenshots/2018-07-05.png)

## Installation

### Requirements

    $ sudo apt install virt-manager curl git
    $ sudo usermod -a -G libvirt $USER # group can also be "libvirtd"
    $ newgrp libvirt
    $ echo user = "\"$USER\"" | sudo tee -a /etc/libvirt/qemu.conf
    $ sudo systemctl restart libvirtd

### Install Nix package manager

    $ sudo mkdir -m 0755 /nix && sudo chown $USER /nix
    $ curl https://nixos.org/nix/install | sh
    $ nix-channel --add https://nixos.org/channels/nixos-19.09 nixpkgs
    $ . ~/.nix-profile/etc/profile.d/nix.sh

### Install appvm

    $ nix-env -if https://code.dumpstack.io/tools/appvm/archive/master.tar.gz

## Usage

### Search for applications

    $ appvm search chromium

### Run application

    $ appvm start chromium
    $ # ... long wait for first time, because we need to collect a lot of packages

### Synchronize remote repos for applications

    $ appvm sync

You can customize local settings in **~/.config/appvm/nix/local.nix**.

Default hotkey to release cursor: ctrl+alt.

### Shared directory

    $ ls appvm/chromium
    foo.tar.gz
    bar.tar.gz

### Close VM

    $ appvm stop chromium

### Automatic ballooning

Add this command:

    $ appvm autoballoon

to crontab like that:

    $ crontab -l
    * * * * * /home/user/dev/go/bin/appvm autoballoon
