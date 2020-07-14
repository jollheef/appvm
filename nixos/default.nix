params@{ config, lib, pkgs, ... }:
let
  cfg = config.virtualisation.appvm;
  appvm = import ../. params;
in with lib; {

  options = {
    virtualisation.appvm = {
      enable = mkOption {
        type = types.bool;
        default = false;
        description = ''
          This enables AppVMs and related virtualisation settings.
        '';
      };
      user = mkOption {
        type = types.str;
        description = ''
          AppVM user login. Currenly only AppVMs are supported for a single user only.
        '';
      };
    };

  };

  config = mkIf cfg.enable {
    virtualisation.libvirtd = {
      enable = true;
      qemuVerbatimConfig = ''
        namespaces = []
        user = "${cfg.user}"
        group = "users"
        remember_owner = 0
      '';
    };

    users.users."${cfg.user}" = {
      packages = [ appvm ];
      extraGroups = [ "libvirtd" ];
    };

  };

}
