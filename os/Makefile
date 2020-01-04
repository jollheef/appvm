test: vm cleanup

vm:
	ln -sf vm.nix target.nix
	nix-build '<nixpkgs/nixos>' -A vm -I nixos-config=configuration.nix
	@echo "Use Ctrl-Alt-Q to close VM" | grep --color=always '.*'
	./result/bin/run-nixos-vm -cpu host

cleanup:
	rm -f nixos.qcow2
	unlink result

iso:
	@echo "Not yet available. Use \`make live-iso\`." | grep --color=always '.*'

live-iso:
	ln -sf live-iso.nix target.nix
	nix-build '<nixpkgs/nixos>' -A config.system.build.isoImage -I nixos-config=configuration.nix
	cp result/iso/* appvm.iso
	unlink result
