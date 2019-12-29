{ pkgs ? import <nixpkgs> {} }:
with pkgs; mkShell {
  buildInputs = [ go gocode virt-viewer virtmanager ];
}
