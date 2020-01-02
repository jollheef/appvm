package main

import (
	"io/ioutil"
)

// Builtin VMs

type app struct {
	Name string
	Nix  []byte
}

var builtin_chromium_nix = app{
	Name: "chromium",
	Nix: []byte(`
{pkgs, ...}:
let
  application = "${pkgs.chromium}/bin/chromium";
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

  programs.chromium = {
    enable = true;
    extensions = [
      "cjpalhdlnbpafiamejdnhcphjbkeiagm" # uBlock Origin
      "gcbommkclmclpchllfjekcdonpmejbdp" # HTTPS Everywhere
      "fihnjjcciajhdojfnbdddfaoknhalnja" # I don't care about cookies
    ];
  };

  services.xserver.displayManager.sessionCommands = "${appRunner}/bin/app &";
}
`),
}

func writeBuiltinApps(path string) (err error) {
	for _, f := range []app{
		builtin_chromium_nix,
	} {
		err = ioutil.WriteFile(configDir+"/nix/"+f.Name+".nix", f.Nix, 0644)
		if err != nil {
			return
		}
	}

	return
}
