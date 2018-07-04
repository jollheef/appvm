{pkgs, ...}:
{
  imports = [
    <nixpkgs/nixos/modules/virtualisation/qemu-vm.nix>
    <nix/base.nix>
  ];

  environment.systemPackages = [ pkgs.libreoffice ];
  services.xserver.displayManager.sessionCommands = "while [ 1 ]; do ${pkgs.libreoffice}/bin/soffice; done &";
}
