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
