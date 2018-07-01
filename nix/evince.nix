{pkgs, ...}:
{
  imports = [
    <nixpkgs/nixos/modules/virtualisation/qemu-vm.nix>
    <nix/base.nix>
  ];

  environment.systemPackages = [ pkgs.evince ];
  services.xserver.displayManager.sessionCommands = "while [ 1 ]; do ${pkgs.evince}/bin/evince; done &";
}
