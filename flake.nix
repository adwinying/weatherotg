{
  description = "weatherotg";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/master";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      devShell = with pkgs; mkShell {
        buildInputs = [
          bun
          go_1_22
          air
          templ
          mprocs
        ];

        shellHook = ''
          echo "`${pkgs.go}/bin/go version`"
          echo "bun: v`${pkgs.bun}/bin/bun --version`"
        '';
      };
    });
}
