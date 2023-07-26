# Ignite nix wrapper
{pkgs}: let
  version = "0.27.1";
  sha256sums = {
    "0.27.1" = "N28aJdE0sfeo/DHVj0GF/KUKge+/1X4EHUNfsu4u5u0=";
  };
  tarball = pkgs.fetchurl {
    url = "https://github.com/ignite/cli/releases/download/v${version}/ignite_${version}_darwin_arm64.tar.gz";
    sha256 = sha256sums."${version}";
  };
in
  pkgs.runCommand "ignite-${version}" {} ''
    mkdir -p $out/bin
    tar xf ${tarball}
    mv ignite $out/bin/
  ''
