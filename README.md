[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=R8W2UQPZ5X5JE&source=url)
[![Donate](https://img.shields.io/badge/Donate-BitCoin-green.svg)](https://blockchair.com/bitcoin/address/bc1q23fyuq7kmngrgqgp6yq9hk8a5q460f39m8nv87)

# Nix application VMs: security through virtualization

Simple application VMs (hypervisor-based sandbox) based on Nix package manager.

Uses one **read-only** /nix directory for all appvms. So creating a new appvm (but not first) is just about one minute.

Currently optimized for full screen usage (but remote-viewer has ability to resize window dynamically without change resolution).

![appvm screenshot](screenshots/2018-07-05.png)

## Dependencies

    $ sudo apt install golang virt-manager curl git
    $ sudo usermod -a -G libvirt $USER

    $ echo 'export GOPATH=$HOME/go' >> ~/.bash_profile
    $ echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bash_profile
    $ echo 'source ~/.bash_profile' >> ~/.bashrc
    $ source ~/.bash_profile

You need to **relogin** if you install virt-manager (libvirt) first time.

## Install Nix package manager

    $ sudo mkdir -m 0755 /nix && sudo chown $USER /nix
    $ curl https://nixos.org/nix/install | sh
    $ . ~/.nix-profile/etc/profile.d/nix.sh

## Libvirt from user (required if you need access to shared files)

    $ echo user = "\"$USER\"" | sudo tee -a /etc/libvirt/qemu.conf
    $ sudo systemctl restart libvirtd

## Install appvm tool

    $ go get code.dumpstack.io/tools/appvm

## Update appvm tool

    $ go get -u code.dumpstack.io/tools/appvm

## Run application

    $ appvm start chromium --verbose
    $ # ... long wait for first time, because we need to collect a lot of packages

You can customize local settings in **~/.config/appvm/nix/local.nix**.

Default hotkey to release cursor: ctrl+alt.

## Shared directory

    $ ls appvm/chromium
    foo.tar.gz
    bar.tar.gz

## Close VM

    $ appvm stop chromium

## Automatic ballooning

Add this command:

    $ appvm autoballoon

to crontab like that:

    $ crontab -l
    * * * * * /home/user/dev/go/bin/appvm autoballoon

# App description

    $ cat nix/chromium.nix
    {pkgs, ...}:
    {
      imports = [
        <nixpkgs/nixos/modules/virtualisation/qemu-vm.nix>
        <nix/base.nix>
      ];

      environment.systemPackages = [ pkgs.chromium ];
      services.xserver.displayManager.sessionCommands = "while [ 1 ]; do ${pkgs.chromium}/bin/chromium; done &";
    }

For create new app you should add package name (search at https://nixos.org/nixos/packages.html) and path to binary (typically same as package name).

## Defined applications (pull requests are welcome!)

* chromium
* thunderbird
* tdesktop
* evince
* libreoffice
* wire
* torbrowser
