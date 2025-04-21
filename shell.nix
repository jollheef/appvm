{ pkgs ? import <nixpkgs> {} }:
with pkgs; mkShell {
  buildInputs = [ go virt-viewer virtmanager ];
}
