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
{
  imports = [
    <nixpkgs/nixos/modules/virtualisation/qemu-vm.nix>
    <nix/base.nix>
  ];

  environment.etc."chromium/policies/managed/plugins.json".text = ''
{
    "ExtensionInstallForcelist": [
        // uBlock Origin (https://chrome.google.com/webstore/detail/ublock-origin/cjpalhdlnbpafiamejdnhcphjbkeiagm)
        "cjpalhdlnbpafiamejdnhcphjbkeiagm;https://clients2.google.com/service/update2/crx",
        // HTTPS Everywhere (https://chrome.google.com/webstore/detail/https-everywhere/gcbommkclmclpchllfjekcdonpmejbdp)
        "gcbommkclmclpchllfjekcdonpmejbdp;https://clients2.google.com/service/update2/crx",
    ]
}
  '';

  environment.systemPackages = [ pkgs.chromium ];
  services.xserver.displayManager.sessionCommands = "while [ 1 ]; do ${pkgs.chromium}/bin/chromium; done &";
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
