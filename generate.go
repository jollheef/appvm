package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/digitalocean/go-libvirt"
)

var template = `
{pkgs, ...}:
{
  imports = [
    <nixpkgs/nixos/modules/virtualisation/qemu-vm.nix>
    <nix/base.nix>
  ];

  services.xserver.displayManager.sessionCommands =
    "while [ 1 ]; do ${pkgs.%s}/bin/%s; done &";
}
`

func isPackageExists(name string) bool {
	cmd := exec.Command("nix-env", "-iA", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err == nil
}

func nixPath(name string) (path string, err error) {
	command := exec.Command("nix", "path-info", name)
	bytes, err := command.Output()
	if err != nil {
		return
	}
	path = string(bytes)
	return
}

func generate(l *libvirt.Libvirt, name, bin string) {
	if !isPackageExists(name) {
		log.Println("Package pkgs."+name, "does not exists")
		return
	}

	path, err := nixPath(name)
	if err != nil {
		log.Println("Cannot find nix path")
		return
	}

	path = strings.TrimSpace(path)

	files, err := ioutil.ReadDir(path + "/bin/")
	if err != nil {
		log.Println(err)
		return
	}

	if bin == "" && len(files) != 1 {
		fmt.Println("There's more than one binary in */bin, " +
			"you should specify one of them explicitly")
		fmt.Println("Files in", path+"/bin/:")
		for _, f := range files {
			fmt.Println("\t", f.Name())
		}
		return
	}

	if bin != "" {
		var found bool = false
		for _, f := range files {
			if bin == f.Name() {
				found = true
			}
		}
		if !found {
			log.Println("There's no such file in */bin")
			return
		}
	} else {
		bin = files[0].Name()
	}

	realName := strings.Split(name, ".")[1]
	appFilename := configDir + "/nix/" + realName + ".nix"

	appNixConfig := fmt.Sprintf(template, realName, bin)

	err = ioutil.WriteFile(appFilename, []byte(appNixConfig), 0600)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Print(appNixConfig + "\n")
	log.Println("Configuration file is saved to", appFilename)
}
