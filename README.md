# Nix app vms

Simple application VM's based on Nix package manager.

Designed primarily for full screen usage without guest additions.

By default uses 3840x2160, so you need to change `appvm/nix/base.nix` monitorSection. Autodetection based on host resolution will be done after :)

## Install Nix package manager

    $ su -c 'mkdir -m 0755 /nix && chown user /nix'
    $ curl https://nixos.org/nix/install | sh

## Dependencies

    $ su -c 'USE="spice virtfs" emerge qemu virt-manager'

## Add appvm to PATH

    $ echo 'PATH=$PATH:$HOME/appvm/bin' >> ~/.bashrc

(if you clone appvm to home directory)

## Create VM

    $ $HOME/appvm/appvm.sh chromium

## Run application

    $ appvm.chromium

## Shared directory

    $ ls appvm/share/chromium
    foo.tar.gz
    bar.tar.gz
