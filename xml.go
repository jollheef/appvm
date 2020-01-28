package main

import "fmt"

// You may think that you want to rewrite to proper golang structures.
// Believe me, you shouldn't.

func generateXML(vmName string, online bool,
	vmNixPath, reginfo, img, sharedDir string) string {

	qemuParams := `
          <qemu:commandline>
            <qemu:arg value='-device'/>
            <qemu:arg value='e1000,netdev=net0'/>
            <qemu:arg value='-netdev'/>
            <qemu:arg value='user,id=net0'/>
            <qemu:arg value='-snapshot'/>
          </qemu:commandline>
        `

	if !online {
		qemuParams = `
                  <qemu:commandline>
                    <qemu:arg value='-snapshot'/>
                  </qemu:commandline>
                `
	}

	return fmt.Sprintf(xmlTmpl, vmName, vmNixPath, vmNixPath, vmNixPath,
		reginfo, img, sharedDir, sharedDir, sharedDir, qemuParams)
}

var xmlTmpl = `
<domain type='kvm' xmlns:qemu='http://libvirt.org/schemas/domain/qemu/1.0'>
  <name>%s</name>
  <memory unit='GiB'>2</memory>
  <currentMemory unit='GiB'>1</currentMemory>
  <vcpu>4</vcpu>
  <os>
    <type arch='x86_64'>hvm</type>
    <kernel>%s/kernel</kernel>
    <initrd>%s/initrd</initrd>
    <cmdline>loglevel=4 init=%s/init %s</cmdline>
  </os>
  <features>
    <acpi></acpi>
  </features>
  <clock offset='utc'/>
  <on_poweroff>destroy</on_poweroff>
  <on_reboot>restart</on_reboot>
  <on_crash>destroy</on_crash>
  <devices>
    <!-- Graphical console -->
    <graphics type='spice' autoport='yes'>
      <listen type='address'/>
      <image compression='off'/>
    </graphics>
    <!-- Guest additionals support -->
    <channel type='spicevmc'>
      <target type='virtio' name='com.redhat.spice.0'/>
    </channel>
    <!-- Fake (because -snapshot) writeback image -->
    <disk type='file' device='disk'>
      <driver name='qemu' type='qcow2' cache='writeback' error_policy='report'/>
      <source file='%s'/>
      <target dev='vda' bus='virtio'/>
    </disk>
    <video>
      <model type='qxl' ram='524288' vram='524288' vgamem='262144' heads='1' primary='yes'/>
      <address type='pci' domain='0x0000' bus='0x00' slot='0x02' function='0x0'/>
    </video>
    <!-- filesystems -->
    <filesystem type='mount' accessmode='passthrough'>
      <source dir='/nix/store'/>
      <target dir='store'/>
      <readonly/>
    </filesystem>
    <filesystem type='mount' accessmode='mapped'>
      <source dir='%s'/>
      <target dir='xchg'/> <!-- workaround for nixpkgs/nixos/modules/virtualisation/qemu-vm.nix -->
    </filesystem>
    <filesystem type='mount' accessmode='mapped'>
      <source dir='%s'/>
      <target dir='shared'/> <!-- workaround for nixpkgs/nixos/modules/virtualisation/qemu-vm.nix -->
    </filesystem>
    <filesystem type='mount' accessmode='mapped'>
      <source dir='%s'/>
      <target dir='home'/>
    </filesystem>
  </devices>
  %s
</domain>
`
