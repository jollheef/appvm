package main

var local_nix_template = []byte(`
{
  services.xserver.layout = "us,ru";
  services.xserver.xkbOptions = "ctrl:nocaps,grp:rctrl_toggle";
}
`)
