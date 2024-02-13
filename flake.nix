{
  description = "Fly.io's distributed systems challenges";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlay = self: super: with self; {
          jepsen-maelstrom = stdenv.mkDerivation {
            name = "maelstrom";
            version = "0.2.3";
            src = builtins.fetchurl {
              url = "https://github.com/jepsen-io/maelstrom/releases/download/v0.2.3/maelstrom.tar.bz2";
              sha256 = "sha256:06jnr113nbnyl9yjrgxmxdjq6jifsjdjrwg0ymrx5pxmcsmbc911";
            };
            installPhase = ''
              mkdir -p $out/bin
              cp -r lib $out/bin/lib
              install -m755 -D maelstrom $out/bin/maelstrom
            '';
          };

          dist-sys-shell = stdenv.mkDerivation {
            name = "dist-sys-shell";
            nativeBuildInputs = [
              go
              openjdk
              graphviz
              gnuplot
              jepsen-maelstrom
            ];
          };
        };

        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;

          overlays = [
            overlay
          ];
        };
      in
      {
        devShell = pkgs.dist-sys-shell;
      }
    );
}
