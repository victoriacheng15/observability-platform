{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  packages = with pkgs; [
    go
  ];

  shellHook = ''
    echo "🚀 go:        $(go version)"
  '';
}