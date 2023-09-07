{rustPlatform}:
rustPlatform.buildRustPackage rec {
  pname = "miden";
  version = "pinned";
  src = ../vendor/github.com/qredek/miden-vm;
  cargoLock.lockFile = ../vendor/github.com/qredek/miden-vm/Cargo.lock;
  buildType = "release";
  buildFeatures = ["executable" "concurrent"];
  # nativeBuildInputs = with nixpkgs; [ rustc ];
  doCheck = false;
}
