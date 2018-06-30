{pkgs, ...}:
{
  imports = [
    <nixpkgs/nixos/modules/virtualisation/qemu-vm.nix>
    <nix/base.nix>
  ];

  environment.systemPackages = [ pkgs.thunderbird ];
  services.xserver.displayManager.sessionCommands = "while [ 1 ]; do ${pkgs.thunderbird}/bin/thunderbird; done &";
}
