{
  description = "Fly.io's distributed system challenges";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          
          config.allowUnfree = true;
        };
      in
      {
        devShell = with pkgs; mkShell {
          buildInputs = [
              go
              go-task
              jet
              graphviz
              gnuplot
              maelstrom-clj
          ];
        };
      }
    );
}
