package main

import (
	"fmt"
	"log"
	"os/user"
)

var base_nix = `
{ cmd, displayProtocol }:
{ pkgs, config, lib, ... }:
let
  protocols = {
    waypipe = {
      services.mingetty = {
        autologinUser = "user";
      };
      hardware.opengl.enable = true;
      fonts.enableDefaultFonts = true;
      programs.dconf.enable = true;
      programs.bash.loginShellInit = ''
        host=$(ip route | grep "^default" | cut -f3 -d" ")
        ${pkgs.socat}/bin/socat TCP-CONNECT:$host:2222 UNIX-LISTEN:/tmp/waypipe-server.sock &
        ${pkgs.waypipe}/bin/waypipe server ${pkgs.dbus}/bin/dbus-run-session ${cmd}
      '';
      environment.variables = {
        QT_QPA_PLATFORM = "wayland";
        XDG_SESSION_TYPE = "wayland";
        QT_WAYLAND_DISABLE_WINDOWDECORATION = "1";
        DISPLAY = ":0";
      };
    };
    virt-viewer = {
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
        displayManager.sessionCommands = "${cmd} &";
        windowManager.xmonad.enable = true;
        windowManager.default = "xmonad";
      };

      services.spice-vdagentd.enable = true;

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
    };
  };
  common = {
    users.extraUsers.user = {
      uid = %s;
      isNormalUser = true;
      extraGroups = [ "audio" ];
      createHome = true;
      password = "";
    };

    systemd.services.mount-home-user = {
      description = "Mount /home/user (crutch)";
      serviceConfig = {
        ExecStart =
          "/bin/sh -c '/run/current-system/sw/bin/mount -t 9p -o trans=virtio,version=9p2000.L home /home/user'";
        RemainAfterExit = "yes";
        Type = "oneshot";
        User = "root";
      };
      wantedBy = [ "sysinit.target" ];
    };

    systemd.services."autoballoon" = {
      serviceConfig = { StartLimitBurst = 100; };
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
      wantedBy = [ "timers.target" ];
    };
  };
in
with lib;
{
  imports = [ <nix/local.nix> ];
  config = common // protocols.${displayProtocol};
}
`

func baseNix() []byte {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return []byte(fmt.Sprintf(base_nix, u.Uid))
}
