{config, options, ...}:
{
  config.networking.hostName = "chromium";
  config.fileSystems."/home/user" = {
    device = "apphome";
    fsType = "9p";
    options = [ "trans=virtio" "version=9p2000.L" "cache=loose" ];
    neededForBoot = true;
  };
}
