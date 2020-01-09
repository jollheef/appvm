package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var template = `
{pkgs, ...}:
let
  application = "${pkgs.%s}/bin/%s";
  appRunner = pkgs.writeShellScriptBin "app" ''
    ARGS_FILE=/home/user/.args
    ARGS=$(cat $ARGS_FILE)
    rm $ARGS_FILE

    ${application} $ARGS
    systemctl poweroff
  '';
in {
  imports = [
    <nixpkgs/nixos/modules/virtualisation/qemu-vm.nix>
    <nix/base.nix>
  ];

  services.xserver.displayManager.sessionCommands = "${appRunner}/bin/app &";
}
`

func isPackageExists(channel, name string) bool {
	return nil == exec.Command("nix-build", "<"+channel+">", "-A", name).Run()
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
			if strings.Contains(fields[1], "nixos.org/channels") {
				channel = fields[0]
				return
			}
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

func generate(pkg, bin, vmname string, build bool) (err error) {
	// TODO refactor
	var name, channel string

	if strings.Contains(pkg, ".") {
		channel = strings.Split(pkg, ".")[0]
		name = strings.Join(strings.Split(pkg, ".")[1:], ".")
	} else {
		log.Println("Package name does not contains channel")
		log.Println("Trying to guess")

		channel, err = guessChannel()
		if err != nil {
			log.Println("Cannot guess channel")
			log.Println("Check nix-channel --list")
			log.Println("Will try <nixpkgs>")
			channel = "nixpkgs"
			err = nil
		}

		name = pkg
		log.Println("Use", channel+"."+pkg)
	}

	if !isPackageExists(channel, name) {
		s := "Package " + name + " does not exists"
		err = errors.New(s)
		log.Println(s)
		return
	}

	path, err := nixPath(channel + "." + name)
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
				parts := strings.Split(pkg, ".")
				if f.Name() == parts[len(parts)-1] {
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

	var appFilename string
	if vmname != "" {
		appFilename = configDir + "/nix/" + vmname + ".nix"
	} else {
		appFilename = configDir + "/nix/" + name + ".nix"
	}

	appNixConfig := fmt.Sprintf(template, name, bin)

	err = ioutil.WriteFile(appFilename, []byte(appNixConfig), 0600)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Print(appNixConfig + "\n")
	log.Println("Configuration file is saved to", appFilename)

	if build {
		if vmname != "" {
			_, _, _, err = generateVM(configDir, vmname, true)
		} else {
			_, _, _, err = generateVM(configDir, name, true)
		}

		if err != nil {
			return
		}
	}

	return
}
