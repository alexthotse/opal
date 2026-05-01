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

        peregrineFrontend = pkgs.buildGoModule {
          pname = "peregrine";
          version = "1.0.0";
          src = ./peregrine;

          # Setting to empty string forces Nix to calculate and show the new hash
          # during the first build if it's incorrect.
          # vendorHash = "sha256-n6V8YgR3bSjJqG9Yy25p50yK2263v7K15E9L41V35G8=";
          vendorHash = pkgs.lib.fakeHash;
          
          # Don't check during build to avoid issues with bubbletea tests in nix sandbox
          doCheck = false;

          nativeBuildInputs = [ pkgs.makeWrapper ];
          buildInputs = pkgs.lib.optionals pkgs.stdenv.isLinux [ pkgs.alsa-lib ];

          postInstall = ''
            wrapProgram $out/bin/peregrine_cli \
              --prefix PATH : ${pkgs.lib.makeBinPath (pkgs.lib.optionals pkgs.stdenv.isLinux [ pkgs.alsa-utils ])}
            
            mv $out/bin/peregrine_cli $out/bin/peregrine
          '';
        };

        alpineBase = pkgs.stdenv.mkDerivation {
          name = "alpine-linux-base.tar";

          nativeBuildInputs = [ pkgs.skopeo ];

          outputHashAlgo = "sha256";
          outputHashMode = "flat";
          outputHash = "sha256-p+iV7+tp4QeS7Scuk9H/TdUk6dwfCQNvr/lqO+eDL9M=";

          buildCommand = ''
            export HOME="$TMPDIR"
            export XDG_RUNTIME_DIR="$TMPDIR"
            skopeo copy docker://alpine@sha256:0a4eaa0eecf5f8c050e5bba433f58c052be7587ee8af3e8b3910ef9ab5fbe9f5 docker-archive:$out:alpine:latest
          '';
        };

        peregrineContainer = pkgs.dockerTools.buildImage {
          name = "ghcr.io/alexthotse/peregrine";
          tag = "latest";
          
          fromImage = alpineBase;

          copyToRoot = pkgs.buildEnv {
            name = "peregrine-env";
            paths = [
              peregrineFrontend
              pkgs.bashInteractive
              pkgs.coreutils
            ] ++ pkgs.lib.optionals pkgs.stdenv.isLinux [ pkgs.alsa-lib pkgs.alsa-utils ];
            pathsToLink = [ "/bin" "/share" ];
          };

          config = {
            Cmd = [ "${peregrineFrontend}/bin/peregrine" ];
            Env = [
              "PATH=/bin:${peregrineFrontend}/bin:${pkgs.bashInteractive}/bin:${pkgs.coreutils}/bin"
            ];
          };
        };
      in
      {
        packages = {
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
          ] ++ pkgs.lib.optionals pkgs.stdenv.isLinux [ alsa-lib alsa-utils ];
        };
      }
    );
}
