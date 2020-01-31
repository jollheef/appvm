let
  pkgs = import <nixpkgs> {};
  virt-manager-without-menu = pkgs.virt-viewer.overrideAttrs(x: {
    patches = [
      ./patches/0001-Remove-menu-bar.patch
      ./patches/0002-Do-not-grab-keyboard-mouse.patch
      ./patches/0003-Use-name-of-appvm-applications-as-a-title.patch
    ];
  });
in with pkgs;

buildGoPackage rec {
  pname = "appvm";
  version = "master";

  buildInputs = [ makeWrapper ];

  goPackagePath = "code.dumpstack.io/tools/${pname}";

  src = ./.;

  goDeps = ./deps.nix;

  postFixup = ''
    wrapProgram $bin/bin/appvm \
      --prefix PATH : "${lib.makeBinPath [ nix virt-manager-without-menu socat waypipe ]}"
  '';

  meta = {
    description = "Nix-based app VMs";
    homepage = "https://code.dumpstack.io/tools/${pname}";
    maintainers = [ lib.maintainers.dump_stack ];
    license = lib.licenses.gpl3;
  };
}
