{pkgs, ...}:
{
  imports = [
    <nixpkgs/nixos/modules/virtualisation/qemu-vm.nix>
    <nix/base.nix>
  ];

  # TODO: block all connections outside tor

  environment.systemPackages = [ pkgs.tor-browser-bundle-bin ];
  services.xserver.displayManager.sessionCommands = "while [ 1 ]; do ${pkgs.tor-browser-bundle-bin}/bin/tor-browser; done &";
}
