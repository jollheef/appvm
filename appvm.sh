#!/bin/bash

APPVM_PATH=$(dirname $(realpath $0))
cd ${APPVM_PATH}

if [[ "$1" == "build" && "$2" != "" ]]; then
    rm qemu/qcow/${2}.qcow2
    NIX_PATH=$NIX_PATH:. nix-build '<nixpkgs/nixos>' -A config.system.build.vm -I nixos-config=nix/${2}.nix || exit 1
    NIX_SYSTEM=$(realpath result/system)
    mkdir -p bin
    RAND_HASH=$(head /dev/urandom | md5sum | awk '{ print $1 }')
    VM_BIN_PATH=$(realpath qemu/bin/qemu.${RAND_HASH}.${2})
    sed "s;NIX_SYSTEM_PLACEHOLDER;${NIX_SYSTEM};" qemu/qemu.template > ${VM_BIN_PATH}
    sed -i "s;NAME_PLACEHOLDER;${2};" ${VM_BIN_PATH}
    sed -i "s;NIX_DISK_IMAGE_PLACEHOLDER;${APPVM_PATH}/qemu/qcow2/${2}.qcow2;" ${VM_BIN_PATH}
    RANDOM_PORT=$(/usr/bin/python -c 'import random; print(random.randint(1024,65535))')
    # TODO Check for port collisions
    sed -i "s;PORT_PLACEHOLDER;${RANDOM_PORT};" ${VM_BIN_PATH}
    echo -e "#!/bin/bash\npgrep -f ${RAND_HASH} || {\n\tnohup setsid ${VM_BIN_PATH} >/dev/null 2>&1 &\n\tsleep 1s\n}\nremote-viewer -f spice://127.200.0.1:${RANDOM_PORT}" > bin/appvm.${2}
    chmod +x ${VM_BIN_PATH}
    chmod +x bin/appvm.${2}
    unlink result
elif [[ "$1" == "generate-resolution" && "$2" != "" && "$3" != "" ]]; then
    MONITOR_SIZE="$(xrandr | grep mm | head -n 1 | awk '{ print $(NF-2) " " $(NF) }' | sed 's/mm//g')"
    CVT="$(cvt ${2} ${3} | grep Modeline)"
    echo "{"
    echo "  services.xserver.monitorSection = ''"
    echo "    " ${CVT}
    echo "    " Option '"PreferredMode"' $(echo ${CVT} | awk '{ print $2 }')
    echo "    " DisplaySize ${MONITOR_SIZE} # In millimeters
    echo "  '';"
    echo "}"
else
    echo -e "Usage:\t$0 build APPLICATION"
    echo -e "or:\t$0 generate-resolution X Y"
fi
