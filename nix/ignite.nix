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
    rev = "7d54608e405e9bad36750c937317c760fb004e12";
    sha256 = "sha256-Va4J6QfzFwNry4wNslPELiU3D5DoLLk812H4cMXnYsQ=";
  };
  vendorSha256 = "sha256-F5+G/eWB/RFldaB8d8a9UQW2fd+afUf0q12+ka940s8=";
  subPackages = ["ignite/cmd/ignite"];
}
