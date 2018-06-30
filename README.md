# Nix app vms

Simple application VM's based on Nix package manager.

Designed primarily for full screen usage without guest additions.

It's a proof-of-concept, but you can still use it. Also there is a lot of strange things inside, don't afraid of :)

## Install Nix package manager

    $ su -c 'mkdir -m 0755 /nix && chown user /nix'
    $ curl https://nixos.org/nix/install | sh

## Dependencies

    $ su -c 'USE="spice virtfs" emerge qemu virt-manager'

## Add appvm to PATH

    $ echo 'PATH=$PATH:$HOME/appvm/bin' >> ~/.bashrc

(if you clone appvm to home directory)

## Generate resolution

By default uses 3840x2160. If you need to regenerate `appvm/nix/monitor.nix`:

    $ appvm/appvm.sh generate-resolution 1920 1080 > appvm/nix/monitor.nix

Autodetection is a bash-spaghetti, so you need to check results. BTW it's just a X.org monitor section.

## Create VM

    $ $HOME/appvm/appvm.sh build chromium

## Run application

    $ appvm.chromium

## Shared directory

    $ ls appvm/share/chromium
    foo.tar.gz
    bar.tar.gz
