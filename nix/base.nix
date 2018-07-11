{
  imports = [
    <nix/monitor.nix>
    <nix/local.nix>
  ];

  system.nixos.stateVersion = "18.03";

  services.xserver = {
    enable = true;
    desktopManager.xterm.enable = false;
    displayManager.slim = {
      enable = true;
      defaultUser = "user";
      autoLogin = true;
    };
    windowManager.xmonad.enable = true;
    windowManager.default = "xmonad";
  };

  users.extraUsers.user = {
    isNormalUser = true;
    extraGroups = [ "audio" ];
    createHome = true;
  };

  environment.etc."xmonad.hs".text = ''
import XMonad
main = xmonad defaultConfig
  { workspaces = [ "" ]
  , borderWidth = 0
  }
  '';

  systemd.services.home-user-build-xmonad = {
    description = "Create and xmonad configuration";
    serviceConfig = {
      ConditionFileNotEmpty = "!/home/user/.xmonad/xmonad.hs";
      ExecStart = "/bin/sh -c 'mkdir -p /home/user/.xmonad && cp /etc/xmonad.hs /home/user/.xmonad/xmonad.hs'";
      RemainAfterExit = "yes";
      Type = "oneshot";
      User = "user";
    };
    wantedBy = [ "multi-user.target" ];
  };

  systemd.services.mount-home-user = {
    description = "Mount /home/user (crutch)";
    serviceConfig = {
      ExecStart = "/bin/sh -c '/run/current-system/sw/bin/mount -t 9p -o trans=virtio,version=9p2000.L,uid=1000 home /home/user'";
      RemainAfterExit = "yes";
      Type = "oneshot";
      User = "root";
    };
    wantedBy = [ "sysinit.target" ];
  };
}
