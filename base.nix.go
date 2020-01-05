package main

var base_nix = []byte(`
{pkgs, ...}:
{
  imports = [
    <nix/local.nix>
  ];

  services.xserver = {
    enable = true;
    desktopManager.xterm.enable = false;
    displayManager.lightdm = {
      enable = true;
      autoLogin = {
        enable = true;
        user = "user";
      };
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
  spawn "while [ 1 ]; do ${pkgs.spice-vdagent}/bin/spice-vdagent -x; done &"
  '';

  systemd.services.home-user-build-xmonad = {
    description = "Link xmonad configuration";
    serviceConfig = {
      ConditionFileNotEmpty = "!/home/user/.xmonad/xmonad.hs";
      ExecStart = "/bin/sh -c 'mkdir -p /home/user/.xmonad && ln -s /etc/xmonad.hs /home/user/.xmonad/xmonad.hs'";
      RemainAfterExit = "yes";
      User = "user";
      Restart = "on-failure";
      TimeoutSec = 10;
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

  systemd.user.services."xrandr" = {
    serviceConfig = {
      StartLimitBurst = 100;
    };
    script = "${pkgs.xorg.xrandr}/bin/xrandr --output Virtual-1 --mode $(${pkgs.xorg.xrandr}/bin/xrandr | grep '   ' | head -n 2 | tail -n 1 | ${pkgs.gawk}/bin/awk '{ print $1 }')";
  };

  systemd.user.timers."xrandr" = {
    description = "Auto update resolution crutch";
    timerConfig = {
      OnBootSec = "1s";
      OnUnitInactiveSec = "1s";
      Unit = "xrandr.service";
      AccuracySec = "1us";
    };
    wantedBy = ["timers.target"];
  };

  systemd.services."autoballoon" = {
    serviceConfig = {
      StartLimitBurst = 100;
    };
    script = ''
      ${pkgs.procps}/bin/free -m | grep Mem | \
        ${pkgs.gawk}/bin/awk '{print $2 "-" $4}' | \
        ${pkgs.bc}/bin/bc > /home/user/.memory_used
    '';
  };

  systemd.timers."autoballoon" = {
    description = "Auto update resolution crutch";
    timerConfig = {
      OnBootSec = "1s";
      OnUnitInactiveSec = "1s";
      Unit = "autoballoon.service";
      AccuracySec = "1us";
    };
    wantedBy = ["timers.target"];
  };
}
`)
