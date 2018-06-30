#!/bin/bash

APPVM_PATH=$(dirname $(realpath $0))
cd ${APPVM_PATH}

if [[ "$1" == "build" && "$2" != "" ]]; then
    NIX_PATH=$NIX_PATH:. nix-build '<nixpkgs/nixos>' -A config.system.build.vm -I nixos-config=nix/${2}.nix || exit 1
    NIX_SYSTEM=$(realpath result/system)
    mkdir -p bin
    sed "s;NIX_SYSTEM_PLACEHOLDER;${NIX_SYSTEM};" qemu/qemu.template > bin/appvm.${2}
    sed -i "s;NAME_PLACEHOLDER;${2};" bin/appvm.${2}
    chmod +x bin/appvm.${2}
    unlink result
else
    echo "Usage: $0 build APPLICATION"
fi
