{pkgs, ...}:
{
  imports = [
    <nixpkgs/nixos/modules/installer/cd-dvd/installation-cd-minimal.nix>
    <chromium-config.nix>
    <base.nix>
  ];

  environment.systemPackages = [ pkgs.chromium ];
  services.xserver.displayManager.sessionCommands = "while [ 1 ]; do ${pkgs.chromium}/bin/chromium; done &";
}
