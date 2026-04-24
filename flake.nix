{
  description = "Opal: A terminal-based AI assistant system";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        opalBackend = pkgs.stdenv.mkDerivation {
          pname = "opal-backend";
          version = "1.0.0";
          src = ./opal_backend;

          nativeBuildInputs = with pkgs; [ gleam erlang rebar3 ];

          buildPhase = ''
            export HOME=$TMPDIR
            gleam deps download
            gleam build
          '';

          installPhase = ''
            mkdir -p $out/share/opal-backend
            cp -r ./* $out/share/opal-backend/
          '';
        };

        opalFrontend = pkgs.buildGoModule {
          pname = "opal-frontend";
          version = "1.0.0";
          src = ./opal_frontend;

          vendorHash = "sha256-n6V8YgR3bSjJqG9Yy25p50yK2263v7K15E9L41V35G8="; # You might need to update this hash
          
          # Don't check during build to avoid issues with bubbletea tests in nix sandbox
          doCheck = false;

          nativeBuildInputs = [ pkgs.makeWrapper ];

          postInstall = ''
            wrapProgram $out/bin/opal_frontend \
              --set OPAL_BACKEND_DIR "${opalBackend}/share/opal-backend" \
              --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.gleam pkgs.erlang ]}
            
            mv $out/bin/opal_frontend $out/bin/opal
          '';
        };
      in
      {
        packages = {
          backend = opalBackend;
          frontend = opalFrontend;
          default = opalFrontend;
        };

        apps.default = flake-utils.lib.mkApp {
          drv = opalFrontend;
          exePath = "/bin/opal";
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
