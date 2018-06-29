# build

    $ nix-build '<nixpkgs/nixos>' -A config.system.build.isoImage -I nixos-config=chromium.nix

# run

   $ qemu-system-x86_64 -smp 2 -m 1024 -enable-kvm -sandbox on -cdrom result/iso/nixos-*-linux.iso
