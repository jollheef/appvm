{
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
}
