{
  description = "Peregrine: A terminal-based AI assistant system";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
    in
    flake-utils.lib.eachSystem supportedSystems (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        peregrineBackend = pkgs.stdenv.mkDerivation {
          pname = "peregrine-backend";
          version = "1.0.0";
          src = ./falcon;

          nativeBuildInputs = with pkgs; [ gleam erlang rebar3 ];

          buildPhase = ''
            export HOME=$TMPDIR
            gleam deps download
            gleam build
          '';

          installPhase = ''
            mkdir -p $out/share/peregrine-backend
            cp -r ./* $out/share/peregrine-backend/
          '';
        };

        peregrineFrontend = pkgs.buildGoModule {
          pname = "peregrine-frontend";
          version = "1.0.0";
          src = ./peregrine;

          vendorHash = "sha256-n6V8YgR3bSjJqG9Yy25p50yK2263v7K15E9L41V35G8="; # You might need to update this hash
          
          # Don't check during build to avoid issues with bubbletea tests in nix sandbox
          doCheck = false;

          nativeBuildInputs = [ pkgs.makeWrapper ];

          postInstall = ''
            wrapProgram $out/bin/peregrine \
              --set FALCON_DIR "${peregrineBackend}/share/peregrine-backend" \
              --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.gleam pkgs.erlang ]}
            
            mv $out/bin/peregrine $out/bin/peregrine
          '';
        };
      in
      {
        packages = {
          backend = peregrineBackend;
          frontend = peregrineFrontend;
          default = peregrineFrontend;
        };

        apps.default = flake-utils.lib.mkApp {
          drv = peregrineFrontend;
          exePath = "/bin/peregrine";
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gleam
            erlang
            rebar3
          ];
        };
      }
    );
}
