/**
 * @author Mikhail Klementev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date July 2018
 * @brief appvm launcher
 */

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/jollheef/go-system"
	"github.com/olekukonko/tablewriter"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var xmlTmpl = `
<domain type='kvm' xmlns:qemu='http://libvirt.org/schemas/domain/qemu/1.0'>
  <name>%s</name>
  <memory unit='GiB'>2</memory>
  <currentMemory unit='GiB'>1</currentMemory>
  <vcpu>4</vcpu>
  <os>
    <type arch='x86_64' machine='pc-i440fx-2.12'>hvm</type>
    <kernel>%s/kernel</kernel>
    <initrd>%s/initrd</initrd>
    <cmdline>loglevel=4 init=%s/init %s</cmdline>
  </os>
  <features>
    <acpi/>
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
  <qemu:commandline>
    <qemu:arg value='-device'/>
    <qemu:arg value='e1000,netdev=net0'/>
    <qemu:arg value='-netdev'/>
    <qemu:arg value='user,id=net0'/>
    <qemu:arg value='-snapshot'/>
  </qemu:commandline>
</domain>
`

func generateXML(name, vmNixPath, reginfo, img, sharedDir string) string {
	// TODO: Define XML in go
	return fmt.Sprintf(xmlTmpl, "appvm_"+name, vmNixPath, vmNixPath, vmNixPath,
		reginfo, img, sharedDir, sharedDir, sharedDir)
}

func list(l *libvirt.Libvirt) {
	domains, err := l.Domains()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Started VM:")
	for _, d := range domains {
		if d.Name[0:5] == "appvm" {
			fmt.Println("\t", d.Name[6:])
		}
	}

	fmt.Println("\nAvailable VM:")
	files, err := ioutil.ReadDir(os.Getenv("GOPATH") + "/src/github.com/jollheef/appvm/nix")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.Name() != "base.nix" &&
			f.Name() != "local.nix" && f.Name() != "monitor.nix" &&
			f.Name() != "local.nix.template" && f.Name() != "monitor.nix.template" {
			fmt.Println("\t", f.Name()[0:len(f.Name())-4])
		}
	}
}

func copyFile(from, to string) (err error) {
	source, err := os.Open(from)
	if err != nil {
		return
	}
	defer source.Close()

	destination, err := os.Create(to)
	if err != nil {
		return
	}

	_, err = io.Copy(destination, source)
	if err != nil {
		destination.Close()
		return
	}

	return destination.Close()
}

func prepareTemplates(appvmPath string) (err error) {
	if _, err = os.Stat(appvmPath + "/nix/local.nix"); os.IsNotExist(err) {
		err = copyFile(appvmPath+"/nix/local.nix.template", appvmPath+"/nix/local.nix")
		if err != nil {
			return
		}
	}

	if _, err = os.Stat(appvmPath + "/nix/monitor.nix"); os.IsNotExist(err) {
		err = copyFile(appvmPath+"/nix/monitor.nix.template", appvmPath+"/nix/monitor.nix")
		if err != nil {
			return
		}
	}
	return
}

func start(l *libvirt.Libvirt, name string) {
	// Currently binary-only installation is not supported, because we need *.nix configurations
	gopath := os.Getenv("GOPATH")
	appvmPath := gopath + "/src/github.com/jollheef/appvm"
	err := os.Chdir(appvmPath)
	if err != nil {
		log.Fatal(err)
	}

	// Copy templates
	err := prepareTemplates(name)
	if err != nil {
		log.Fatal(err)
	}

	stdout, stderr, ret, err := system.System("nix-build", "<nixpkgs/nixos>", "-A", "config.system.build.vm",
		"-I", "nixos-config=nix/"+name+".nix", "-I", ".")
	if err != nil {
		log.Fatalln(err, stdout, stderr, ret)
	}

	realpath, err := filepath.EvalSymlinks("result/system")
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Use go regex
	reginfo, _, _, err := system.System("sh", "-c", "cat result/bin/run-nixos-vm | grep -o 'regInfo=.*/registration'")
	if err != nil {
		log.Fatal(err)
	}

	syscall.Unlink("result")

	qcow2 := "/tmp/.appvm.fake.qcow2"
	if _, err := os.Stat(qcow2); os.IsNotExist(err) {
		system.System("qemu-img", "create", "-f", "qcow2", qcow2, "512M")
		err := os.Chmod(qcow2, 0400) // qemu run with -snapshot, we only need it for create /dev/vda
		if err != nil {
			log.Fatal(err)
		}
	}

	sharedDir := fmt.Sprintf(os.Getenv("HOME") + "/appvm/" + name)
	os.MkdirAll(sharedDir, 0700)

	// TODO: Search go libraries for manipulate ACL
	_, _, _, err = system.System("setfacl", "-R", "-m", "u:qemu:rwx", os.Getenv("HOME")+"/appvm/")
	if err != nil {
		log.Fatal(err)
	}

	xml := generateXML(name, realpath, reginfo, qcow2, sharedDir)
	_, err = l.DomainCreateXML(xml, libvirt.DomainStartValidate)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("virt-viewer", "appvm_"+name)
	cmd.Start()
}

func stop(l *libvirt.Libvirt, name string) {
	dom, err := l.DomainLookupByName("appvm_" + name)
	if err != nil {
		if libvirt.IsNotFound(err) {
			log.Println("Appvm not found or already stopped")
			return
		} else {
			log.Fatal(err)
		}
	}
	err = l.DomainShutdown(dom)
	if err != nil {
		log.Fatal(err)
	}
}

func drop(name string) {
	appDataPath := fmt.Sprintf(os.Getenv("HOME") + "/appvm/" + name)
	os.RemoveAll(appDataPath)
}

func autoBalloon(l *libvirt.Libvirt, memoryMin, adjustPercent uint64) {
	domains, err := l.Domains()
	if err != nil {
		log.Fatal(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Application VM", "Used memory", "Current memory", "Max memory", "New memory"})
	for _, d := range domains {
		if d.Name[0:5] == "appvm" {
			name := d.Name[6:]

			memoryUsedRaw, err := ioutil.ReadFile(os.Getenv("HOME") + "/appvm/" + name + "/.memory_used")
			if err != nil {
				log.Fatal(err)
			}
			memoryUsedMiB, err := strconv.Atoi(string(memoryUsedRaw[0 : len(memoryUsedRaw)-1]))
			if err != nil {
				log.Fatal(err)
			}
			memoryUsed := memoryUsedMiB * 1024

			_, memoryMax, memoryCurrent, _, _, err := l.DomainGetInfo(d)
			if err != nil {
				log.Fatal(err)
			}

			memoryNew := uint64(float64(memoryUsed) * (1 + float64(adjustPercent)/100))

			if memoryNew > memoryMax {
				memoryNew = memoryMax - 1
			}

			if memoryNew < memoryMin {
				memoryNew = memoryMin
			}

			err = l.DomainSetMemory(d, memoryNew)
			if err != nil {
				log.Fatal(err)
			}

			table.Append([]string{name,
				fmt.Sprintf("%d", memoryUsed),
				fmt.Sprintf("%d", memoryCurrent),
				fmt.Sprintf("%d", memoryMax),
				fmt.Sprintf("%d", memoryNew)})
		}
	}
	table.Render()
}

func main() {
	c, err := net.DialTimeout("unix", "/var/run/libvirt/libvirt-sock", time.Second)
	if err != nil {
		log.Fatal(err)
	}

	l := libvirt.New(c)
	if err := l.Connect(); err != nil {
		log.Fatal(err)
	}
	defer l.Disconnect()

	kingpin.Command("list", "List applications")
	autoballonCommand := kingpin.Command("autoballoon", "Automatically adjust/reduce app vm memory")
	minMemory := autoballonCommand.Flag("min-memory", "Set minimal memory (megabytes)").Default("1024").Uint64()
	adjustPercent := autoballonCommand.Flag("adj-memory", "Adjust memory ammount (percents)").Default("20").Uint64()
	startName := kingpin.Command("start", "Start application").Arg("name", "Application name").Required().String()
	stopName := kingpin.Command("stop", "Stop application").Arg("name", "Application name").Required().String()
	dropName := kingpin.Command("drop", "Remove application data").Arg("name", "Application name").Required().String()

	switch kingpin.Parse() {
	case "list":
		list(l)
	case "start":
		start(l, *startName)
	case "stop":
		stop(l, *stopName)
	case "drop":
		drop(*dropName)
	case "autoballoon":
		autoBalloon(l, *minMemory*1024, *adjustPercent)
	}
}
