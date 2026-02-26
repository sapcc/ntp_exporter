# SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company
# SPDX-License-Identifier: Apache-2.0

{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell {
  nativeBuildInputs = [
    addlicense
    go-licence-detector
    go_1_26
    golangci-lint
    goreleaser
    gotools # goimports
    renovate
    reuse
    syft
    typos
    # keep this line if you use bash
    bashInteractive
  ];
}
