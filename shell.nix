{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell {
  nativeBuildInputs = [
    go-licence-detector
    go_1_23
    golangci-lint
    goreleaser
    gotools # goimports

    # keep this line if you use bash
    bashInteractive
  ];
}
