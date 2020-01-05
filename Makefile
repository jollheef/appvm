go:
	go build

nix:
	nix-build -E '((import <nixpkgs> {}).callPackage (import ./default.nix) { })' --option sandbox true --no-out-link
