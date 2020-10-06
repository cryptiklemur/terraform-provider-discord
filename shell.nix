{ pkgs ? import ./default.nix { } }:
with pkgs;
mkShell {
  buildInputs = [
    bash
    curl
    direnv
    dos2unix
    git
    git-lfs
    github-release
    glibcLocales
    gnumake
    go
    nixfmt
    nix-prefetch-git
    which
    wget
  ];
}
