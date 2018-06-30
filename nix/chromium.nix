{pkgs, ...}:
{
  imports = [
    <nixpkgs/nixos/modules/virtualisation/qemu-vm.nix>
    <nix/base.nix>
  ];

  environment.systemPackages = [ pkgs.chromium ];
  services.xserver.displayManager.sessionCommands = "while [ 1 ]; do ${pkgs.chromium}/bin/chromium; done &";
}
