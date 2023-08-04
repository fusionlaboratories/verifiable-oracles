{
  pkgs,
  lib,
  ...
}: let
  ignite = pkgs.callPackage ./nix/ignite.nix {};
in {
  # https://devenv.sh/basics/
  env.GREET = "devenv";

  # https://devenv.sh/packages/
  packages = [ignite pkgs.nodePackages.ganache];

  # https://devenv.sh/scripts/
  # scripts.hello.exec = "echo hello from $GREET";

  enterShell = ''
    export GOPATH="$(go env GOPATH)"
    echo GOPATH=$GOPATH
    ignite version
    ganache --version
  '';

  # https://devenv.sh/languages/
  # languages.nix.enable = true;
  languages.go.enable = true;
  env.GOPATH = lib.mkForce "";

  # https://devenv.sh/pre-commit-hooks/
  # pre-commit.hooks.shellcheck.enable = true;
  pre-commit.hooks.alejandra.enable = true;
  pre-commit.hooks.gofmt.enable = true;

  # https://devenv.sh/processes/
  # processes.ping.exec = "ping example.com";

  # See full reference at https://devenv.sh/reference/options/
}
