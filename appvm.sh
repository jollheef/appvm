#!/bin/bash

APPVM_PATH=$(dirname $(realpath $0))
cd ${APPVM_PATH}

if [[ "$1" == "build" && "$2" != "" ]]; then
    NIX_PATH=$NIX_PATH:. nix-build '<nixpkgs/nixos>' -A config.system.build.vm -I nixos-config=nix/${2}.nix || exit 1
    NIX_SYSTEM=$(realpath result/system)
    mkdir -p bin
    sed "s;NIX_SYSTEM_PLACEHOLDER;${NIX_SYSTEM};" qemu/qemu.template > bin/appvm.${2}
    sed -i "s;NAME_PLACEHOLDER;${2};" bin/appvm.${2}
    sed -i "s;NIX_DISK_IMAGE_PLACEHOLDER;${APPVM_PATH}/qemu/qcow2/${2}.qcow2;" bin/appvm.${2}
    RANDOM_PORT=$(/usr/bin/python -c 'import random; print(random.randint(1024,65535))')
    # TODO Check for port collisions
    sed -i "s;PORT_PLACEHOLDER;${RANDOM_PORT};" bin/appvm.${2}
    echo -e "#!/bin/bash\nremote-viewer -f spice://127.200.0.1:${RANDOM_PORT}" > bin/appgui.${2}
    chmod +x bin/app{vm,gui}.${2}
    unlink result
else
    echo "Usage: $0 build APPLICATION"
fi
