# Ignite nix wrapper, currently only supports arm64
{
  pkgs,
  stdenv,
}: let
  version = "0.27.1";
  system =
    if stdenv.isDarwin
    then "darwin"
    else "linux";
  sha256sums = {
    darwin = {
      "0.27.1" = "N28aJdE0sfeo/DHVj0GF/KUKge+/1X4EHUNfsu4u5u0=";
    };
    linux = {
      "0.27.1" = "/kujtVVpraAtZMJU/pWb9pQG6Bz6/1sTq1Yz6wxgckU=";
    };
  };
  tarball = pkgs.fetchurl {
    url = "https://github.com/ignite/cli/releases/download/v${version}/ignite_${version}_${system}_arm64.tar.gz";
    sha256 = sha256sums."${system}"."${version}";
  };
in
  pkgs.runCommand "ignite-${version}" {} ''
    mkdir -p $out/bin
    tar xf ${tarball}
    mv ignite $out/bin/
  ''
