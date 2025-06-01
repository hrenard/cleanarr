{ pkgs, lib, ... }:
pkgs.buildGoModule {
  name = "cleanarr";
  src = lib.fileset.toSource {
    root = ./.;
    fileset = lib.fileset.intersection ./. (
      lib.fileset.unions [
        ./go.mod
        ./go.sum
        ./main.go
        ./internal
      ]
    );
  };
  vendorHash = "sha256-y1eIpKBqGfLfYw3eds+TPUbN5/PjDTtY1P334YaiBwg=";
  meta = {
    mainProgram = "cleanarr";
    description = "Utility tasked to automatically clean radarr and sonarr files over time";
    homepage = "https://github.com/hrenard/cleanarr";
    license = lib.licenses.gpl3;
    platforms = lib.platforms.linux;
    # maintainers = with lib.maintainers; [ hougo ];
  };
}
