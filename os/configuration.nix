{ config, pkgs, lib, ... }:

let
  appvm = (pkgs.buildGoPackage {
    # TODO ../default.nix
    name = "appvm";
    goPackagePath = "code.dumpstack.io/tools/appvm";
    goDeps = ../deps.nix;
    src = builtins.fetchGit {
      url = "https://code.dumpstack.io/tools/appvm.git";
      ref = "master";
    };
    buildInputs = [ pkgs.makeWrapper ];
    postFixup = ''
      wrapProgram $bin/bin/appvm \
        --prefix PATH : "${lib.makeBinPath [ pkgs.nix pkgs.virt-viewer ]}"
    '';
  });
in {
  imports = [
    ./target.nix
    #./hardware-configuration.nix
  ];

  time.timeZone = "UTC";

  boot.loader.systemd-boot.enable = true;

  # You can not use networking.networkmanager with networking.wireless
  networking.wireless.enable = false;

  systemd.services."init-nix-channels" = {
    enable = true;
    serviceConfig = {
      ExecStartPre = "${pkgs.su}/bin/su root -c '${pkgs.nix}/bin/nix-channel --update'";
      ExecStart = "/bin/sh";
      Restart = "on-failure";
      RestartSec = "5";
      TimeoutSec = "120";
    };
  };

  systemd.timers."init-nix-channels" = {
    timerConfig.OnBootSec = "30s";
    timerConfig.Unit = "init-nix-channels.service";
    wantedBy = ["timers.target"];
  };

  users.users.user = {
    isNormalUser = true;
    extraGroups = [ "audio" "libvirtd" ];
  };

  virtualisation.libvirtd = {
    enable = true;
    qemuVerbatimConfig = ''
      namespaces = []
      user = "user"
      group = "users"
    '';
  };

  # TODO run ${appvm}/bin/appvm autoballoon each second

  environment.systemPackages = with pkgs; [
    appvm virtmanager chromium
    # Cache packages required for application VMs
    xmonad-with-packages spice-vdagent bc qemu_test slim
  ];

  services.xserver.enable = true;
  services.xserver.displayManager.gdm = {
    enable = true;
    wayland = false;            # FIXME
    autoLogin = {
      enable = true;
      user = "user";
    };
  };

  services.xserver.desktopManager.gnome3.enable = true;
}
