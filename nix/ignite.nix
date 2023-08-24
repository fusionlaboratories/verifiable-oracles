# Ignite nix wrapper, currently only supports arm64
{
  pkgs,
  lib,
}:
pkgs.buildGoModule rec {
  pname = "ignite";
  version = src.rev;
  src = pkgs.fetchFromGitHub {
    owner = "ignite";
    repo = "cli";
    rev = "6cb985b5523b34d29dd99c05b15c72c10ea838a3";
    sha256 = "sha256-X07b6/+OIv37IgYDAhkNU+16xiL9xNIzP049SmU+r28=";
  };
  vendorSha256 = "sha256-F5+G/eWB/RFldaB8d8a9UQW2fd+afUf0q12+ka940s8=";
  subPackages = ["ignite/cmd/ignite"];
}
