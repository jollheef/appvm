{pkgs, ...}:
{
  imports = [
    <nixpkgs/nixos/modules/virtualisation/qemu-vm.nix>
    <nix/base.nix>
  ];

  environment.systemPackages = [ pkgs.wire-desktop ];
  services.xserver.displayManager.sessionCommands = "while [ 1 ]; do ${pkgs.wire-desktop}/bin/wire-desktop; done &";
}
