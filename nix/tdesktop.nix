{pkgs, ...}:
{
  imports = [
    <nixpkgs/nixos/modules/virtualisation/qemu-vm.nix>
    <nix/base.nix>
  ];

  environment.systemPackages = [ pkgs.tdesktop ];
  services.xserver.displayManager.sessionCommands = "while [ 1 ]; do ${pkgs.tdesktop}/bin/telegram-desktop; done &";
}
