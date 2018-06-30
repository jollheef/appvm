# Nix app vms

## Install Nix package manager

    $ su -c 'mkdir -m 0755 /nix && chown user /nix'
    $ curl https://nixos.org/nix/install | sh

## Create VM

    $ NIX_PATH=$NIX_PATH:. nix-build '<nixpkgs/nixos>' -A config.system.build.vm -I nixos-config=chromium.nix
    $ ln -s
    these derivations will be built:
      /nix/store/alsxfss45f61015qk0fi147iidl3hj7h-system-path.drv
      /nix/store/y4r9v7x4wh2pa1slrv0jw97bzax19ssv-dbus-1.drv
    ...
    these paths will be fetched (80.31 MiB download, 341.86 MiB unpacked):
      /nix/store/2bqrxp8j4ax2ka5car249bb8jdnh3rvm-adwaita-icon-theme-3.28.0
      /nix/store/rrxkl46g4dc8140ykh9i0clswvfqmz1g-chromium-67.0.3396.87-sandbox
    ...
    building '/nix/store/sbxmqqi3wmpf9n79a0mncrvb302xwh4n-nixos-vm.drv'...
    /nix/store/7zgnq9xw9i7ybvjg0d3h2ank5g5nya5p-nixos-vm

## Run application

For last built application (result is a symbolic link):

    $ ./result/bin/run-chromium-vm

Or by full path:

    $ /nix/store/*/bin/run-chromium-vm

