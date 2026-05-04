{
  description = "Peregrine: A terminal-based AI assistant system";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gastown.url = "github:gastownhall/gastown";
    beads.url = "github:gastownhall/beads";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    gastown,
    beads,
  }:
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

        gtPkg = gastown.packages.${system}.default;
        bdPkg = beads.packages.${system}.default;

        peregrineFrontend = pkgs.buildGoModule {
          pname = "peregrine";
          version = "1.0.0";
          src = ./peregrine;

          # Setting to empty string forces Nix to calculate and show the new hash
          # during the first build if it's incorrect.
          # vendorHash = "sha256-n6V8YgR3bSjJqG9Yy25p50yK2263v7K15E9L41V35G8=";
          vendorHash = "sha256-IdMpy1B4BbilibgbTbO4xoM162l69mF69YGCQHLS5OE=";
          
          # Don't check during build to avoid issues with bubbletea tests in nix sandbox
          doCheck = false;

          nativeBuildInputs = [ pkgs.makeWrapper ];
          buildInputs = pkgs.lib.optionals pkgs.stdenv.isLinux [ pkgs.alsa-lib ];

          postInstall = ''
            wrapProgram $out/bin/peregrine \
              --prefix PATH : ${pkgs.lib.makeBinPath (pkgs.lib.optionals pkgs.stdenv.isLinux [ pkgs.alsa-utils ])}
          '';
        };

        harness = pkgs.writeShellApplication {
          name = "opal-harness";
          runtimeInputs = [
            bdPkg
            gtPkg
            pkgs.bash
            pkgs.coreutils
            pkgs.git
          ];
          text = builtins.readFile ./harness/gt.sh;
        };

        bundle = pkgs.buildEnv {
          name = "opal";
          paths = [
            bdPkg
            gtPkg
            harness
            peregrineFrontend
          ];
          pathsToLink = [ "/bin" "/share" ];
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
          bd = bdPkg;
          gt = gtPkg;
          harness = harness;
          bundle = bundle;
          frontend = peregrineFrontend;
          container = peregrineContainer;
          default = bundle;
        };

        apps = {
          peregrine = flake-utils.lib.mkApp {
            drv = peregrineFrontend;
            exePath = "/bin/peregrine";
          };
          harness = flake-utils.lib.mkApp { drv = harness; };
          default = self.apps.${system}.peregrine;
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
            gtPkg
            bdPkg
          ] ++ pkgs.lib.optionals pkgs.stdenv.isLinux [ alsa-lib alsa-utils ];
        };
      }
    );
}
