# Nix app vms

Simple application VM's based on Nix package manager.

Uses one **read-only** /nix directory for all appvms. So creating a new appvm (but not first) is just about one minute.

Currently optimized for full screen usage (but remote-viewer has ability to resize window dynamically without change resolution) without guest additions.

![appvm screenshot](screenshots/2018-07-05.png)

## Install Nix package manager

    $ su -c 'mkdir -m 0755 /nix && chown user /nix'
    $ curl https://nixos.org/nix/install | sh

## Dependencies

    $ su -c 'USE="spice virtfs" emerge qemu virt-manager'

## Libvirt from user (required if you need access to shared files)

    $ echo user = "\"$USER\"" | sudo tee -a /etc/libvirt/qemu.conf

## Install appvm tool

    $ go get github.com/jollheef/appvm

## Generate resolution

By default uses 3840x2160. If you need to regenerate `appvm/nix/monitor.nix`:

    $ $GOPATH/github.com/jollheef/appvm/generate-resolution.sh 1920 1080 > $GOPATH/github.com/jollheef/appvm/nix/monitor.nix

Autodetection is a bash-spaghetti, so you need to check results. BTW it's just a X.org monitor section.

## Run application

($GOPATH/bin must be in $PATH)

    $ appvm start chromium

You can customize local settings in `$GOPATH/github.com/jollheef/appvm/nix/local.nix`.

Default hotkey to release cursor: ctrl+alt.

## Shared directory

    $ ls appvm/chromium
    foo.tar.gz
    bar.tar.gz

## Close VM

    $ appvm stop chromium

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
