{pkgs, ...}:
{
  imports = [
    <nix/monitor.nix>
    <nix/local.nix>
  ];

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

  services.spice-vdagentd.enable = true;

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
  , startupHook = startup
  }

startup :: X ()
startup = do
  spawn "spice-vdagent"
  '';

  environment.systemPackages = [ pkgs.bc ];
  services.cron = {
    enable = true;
    systemCronJobs = [
      "* * * * *      root    free -m | grep Mem | awk '{print $2 \"-\" $4}' | bc > /home/user/.memory_used"
    ];
  };

  systemd.services.home-user-build-xmonad = {
    description = "Create and xmonad configuration";
    serviceConfig = {
      ConditionFileNotEmpty = "!/home/user/.xmonad/xmonad.hs";
      ExecStart = "/bin/sh -c 'mkdir -p /home/user/.xmonad && ln -s /etc/xmonad.hs /home/user/.xmonad/xmonad.hs'";
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
