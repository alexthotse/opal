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

        falconBackend = pkgs.stdenv.mkDerivation {
          pname = "falcon";
          version = "1.0.0";
          src = ./falcon;

          nativeBuildInputs = with pkgs; [ gleam erlang rebar3 ];

          buildPhase = ''
            export HOME=$TMPDIR
            gleam deps download
            gleam build
          '';

          installPhase = ''
            mkdir -p $out/share/falcon
            cp -r ./* $out/share/falcon/
          '';
        };

        peregrineFrontend = pkgs.buildGoModule {
          pname = "peregrine";
          version = "1.0.0";
          src = ./peregrine;

          # Setting to empty string forces Nix to calculate and show the new hash
          # during the first build if it's incorrect.
          # vendorHash = "sha256-n6V8YgR3bSjJqG9Yy25p50yK2263v7K15E9L41V35G8=";
          vendorHash = null;
          
          # Don't check during build to avoid issues with bubbletea tests in nix sandbox
          doCheck = false;

          nativeBuildInputs = [ pkgs.makeWrapper ];
          buildInputs = [ pkgs.alsa-lib ];

          postInstall = ''
            wrapProgram $out/bin/peregrine_cli \
              --set FALCON_DIR "${falconBackend}/share/falcon" \
              --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.gleam pkgs.erlang ]}
            
            mv $out/bin/peregrine_cli $out/bin/peregrine
          '';
        };

        # OCI (Docker) container image using Alpine Linux as the base
        peregrineContainer = pkgs.dockerTools.buildImage {
          name = "ghcr.io/alexthotse/peregrine";
          tag = "latest";
          
          # Pull a minimal Alpine Linux image as the base
          fromImage = pkgs.dockerTools.pullImage {
            imageName = "alpine";
            imageDigest = "sha256:0a4eaa0eecf5f8c050e5bba433f58c052be7587ee8af3e8b3910ef9ab5fbe9f5"; # alpine:3.21.3
            sha256 = "sha256-p+iV7+tp4QeS7Scuk9H/TdUk6dwfCQNvr/lqO+eDL9M=";
            finalImageName = "alpine";
            finalImageTag = "latest";
          };

          copyToRoot = pkgs.buildEnv {
            name = "peregrine-env";
            paths = [
              peregrineFrontend
              pkgs.bashInteractive
              pkgs.coreutils
              pkgs.alsa-lib
            ];
            pathsToLink = [ "/bin" "/share" ];
          };

          config = {
            Cmd = [ "${peregrineFrontend}/bin/peregrine" ];
            Env = [
              "PATH=${peregrineFrontend}/bin:${pkgs.bashInteractive}/bin:${pkgs.coreutils}/bin"
            ];
          };
        };
      in
      {
        packages = {
          backend = falconBackend;
          frontend = peregrineFrontend;
          container = peregrineContainer;
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
            go-task
            bazelisk
            buf
            protobuf
          ];
        };
      }
    );
}
