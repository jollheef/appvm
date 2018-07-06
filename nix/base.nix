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
      ExecStart = "/bin/sh -c 'mkdir /home/user/.xmonad && cp /etc/xmonad.hs /home/user/.xmonad/xmonad.hs'";
      RemainAfterExit = "yes";
      Type = "oneshot";
      User = "user";
    };
    wantedBy = [ "multi-user.target" ];
  };
}
