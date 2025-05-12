# SPDX-FileCopyrightText: 2024 SAP SE
# SPDX-License-Identifier: Apache-2.0

{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell {
  nativeBuildInputs = [
    addlicense
    go-licence-detector
    go_1_24
    golangci-lint
    goreleaser
    gotools # goimports
    # keep this line if you use bash
    bashInteractive
  ];
}
