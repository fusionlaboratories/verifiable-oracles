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
  packages = [ignite pkgs.nodePackages.ganache pkgs.gotestsum];

  # https://devenv.sh/scripts/
  # devenv splash screen, not run in CI
  scripts.devenv-splash.exec = ''
    echo GOPATH=$GOPATH
    ignite version
    ganache --version
  '';
  # This scripts starts ganache and writes ganache.env file, which contains all the environment variables
  # TODO Use $CI_JOB_ID to add a suffic to ganache.env
  scripts.ganache-host.exec = "ganache instances list | grep $1 | sed -e 's/│//g' | awk '{ print $5 }'";
  scripts.ganache-start.exec = ''
    GANACHE=$(ganache --wallet.deterministic=true --detach)

    # env file
    > ganache.env
    echo "GANACHE=\"$GANACHE\"" | tee -a ganache.env
  '';
  scripts.ganache-stop.exec = "ganache instances stop $1";

  enterShell = ''
    export GOPATH="$(go env GOPATH)"
    if [ "$GITLAB_CI" != "true" ]; then
      devenv-splash
    fi
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
  processes.ganache.exec = "ganache --wallet.deterministic=true";

  # See full reference at https://devenv.sh/reference/options/
}
