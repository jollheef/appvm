package main

import (
	"errors"
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

func guessChannel() (channel string, err error) {
	command := exec.Command("nix-channel", "--list")
	bytes, err := command.Output()
	if err != nil {
		return
	}
	channels := strings.Split(string(bytes), "\n")
	for _, line := range channels {
		fields := strings.Fields(line)
		if len(fields) == 2 {
			channel = fields[0]
			return
		}
	}

	err = errors.New("No channel found")
	return
}

func filterDotfiles(files []os.FileInfo) (notHiddenFiles []os.FileInfo) {
	for _, f := range files {
		if !strings.HasPrefix(f.Name(), ".") {
			notHiddenFiles = append(notHiddenFiles, f)
		}
	}
	return
}

func generate(l *libvirt.Libvirt, pkg, bin, vmname string) {
	var name string

	if strings.Contains(pkg, ".") {
		name = pkg
	} else {
		log.Println("Package name does not contains channel")
		log.Println("Trying to guess")
		channel, err := guessChannel()
		if err != nil {
			log.Println("Cannot guess channel")
			log.Println("Check nix-channel --list")
			return
		}

		name = channel + "." + pkg
		log.Println("Use", name)
	}

	if !isPackageExists(name) {
		log.Println("Package", name, "does not exists")
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
		fmt.Println("There's more than one binary in */bin")
		fmt.Println("Files in", path+"/bin/:")
		for _, f := range files {
			fmt.Println("\t", f.Name())
		}

		log.Println("Trying to guess binary")
		var found bool = false

		notHiddenFiles := filterDotfiles(files)
		if len(notHiddenFiles) == 1 {
			log.Println("Use", notHiddenFiles[0].Name())
			bin = notHiddenFiles[0].Name()
			found = true
		}

		if !found {
			for _, f := range files {
				if f.Name() == pkg {
					log.Println("Use", f.Name())
					bin = f.Name()
					found = true
				}
			}
		}

		if !found {
			log.Println("Cannot guess in */bin, " +
				"you should specify one of them explicitly")
			return
		}
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

	var appFilename string
	if vmname != "" {
		appFilename = configDir + "/nix/" + vmname + ".nix"
	} else {
		appFilename = configDir + "/nix/" + realName + ".nix"
	}

	appNixConfig := fmt.Sprintf(template, realName, bin)

	err = ioutil.WriteFile(appFilename, []byte(appNixConfig), 0600)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Print(appNixConfig + "\n")
	log.Println("Configuration file is saved to", appFilename)
}
