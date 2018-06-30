{
  imports = [
    <nix/monitor.nix>
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
    monitorSection = ''
      Modeline "3840x2160_60.00"  712.75  3840 4160 4576 5312  2160 2163 2168 2237 -hsync +vsync
      Option "PreferredMode" "3840x2160_60.00"
      DisplaySize 610 350    # In millimeters
    '';
  };

  users.extraUsers.user = {
    isNormalUser = true;
    extraGroups = [ "audio" ];
    createHome = true;
  };
}
