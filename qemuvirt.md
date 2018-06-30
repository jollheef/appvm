# /usr/bin/qemu-system-x86_64 

## need

    -sandbox on,obsolete=deny,elevateprivileges=deny,spawn=deny,resourcecontrol=deny
    -m 1024
    -realtime mlock=off
    -smp 1,sockets=1,cores=1,threads=1
    -drive file=/home/user/appvm/result/iso/nixos-18.09pre144732.2171fc4d55b-x86_64-linux.iso,format=raw,if=none,id=drive-ide0-0-0,readonly=on
    -chardev spicevmc,id=charchannel0,name=vdagent
    -device virtserialport,bus=virtio-serial0.0,nr=1,chardev=charchannel0,id=channel0,name=com.redhat.spice.0
    -spice port=5900,addr=127.0.0.1,disable-ticketing,image-compression=off,seamless-migration=on
    -device qxl-vga,id=video0,ram_size=67108864,vram_size=67108864,vram64_size_mb=0,vgamem_mb=16,max_outputs=1,bus=pci.0,addr=0x2
    -device intel-hda,id=sound0,bus=pci.0,addr=0x3

## no need :)

    -msg timestamp=on
    -name guest=generic,debug-threads=on 
    -S 
    -object secret,id=masterKey0,format=raw,file=/var/lib/libvirt/qemu/domain-1-generic/master-key.aes
    -machine pc-i440fx-2.12,accel=kvm,usb=off,vmport=off,dump-guest-core=off
    -cpu Skylake-Client-IBRS
    -device usb-redir,chardev=charredir1,id=redir1,bus=usb.0,port=2
    -device virtio-balloon-pci,id=balloon0,bus=pci.0,addr=0x6
    -device hda-duplex,id=sound0-codec0,bus=sound0.0,cad=0
    -chardev spicevmc,id=charredir0,name=usbredir
    -device usb-redir,chardev=charredir0,id=redir0,bus=usb.0,port=1
    -chardev spicevmc,id=charredir1,name=usbredir
    -device isa-serial,chardev=charserial0,id=serial0
    -device ide-cd,bus=ide.0,unit=0,drive=drive-ide0-0-0,id=ide0-0-0,bootindex=1
    -chardev pty,id=charserial0
    -boot strict=on
    -device ich9-usb-ehci1,id=usb,bus=pci.0,addr=0x4.0x7
    -device ich9-usb-uhci1,masterbus=usb.0,firstport=0,bus=pci.0,multifunction=on,addr=0x4
    -device ich9-usb-uhci2,masterbus=usb.0,firstport=2,bus=pci.0,addr=0x4.0x1
    -device ich9-usb-uhci3,masterbus=usb.0,firstport=4,bus=pci.0,addr=0x4.0x2
    -device virtio-serial-pci,id=virtio-serial0,bus=pci.0,addr=0x5
    -uuid cf00561e-ddc2-4be4-a91a-08b0c28f3f75
    -no-user-config
    -nodefaults
    -chardev socket,id=charmonitor,path=/var/lib/libvirt/qemu/domain-1-generic/monitor.sock,server,nowait
    -mon chardev=charmonitor,id=monitor,mode=control
    -rtc base=utc,driftfix=slew
    -global kvm-pit.lost_tick_policy=delay
    -no-hpet
    -no-shutdown
    -global PIIX4_PM.disable_s3=1
    -global PIIX4_PM.disable_s4=1
