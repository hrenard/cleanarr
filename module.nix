{
  config,
  lib,
  pkgs,
  ...
}:
let
  inherit (lib) types;
  cleanarr = pkgs.callPackage ./package.nix { };
  cfg = config.services.cleanarr;
  servarrSettingType = types.submodule {
    options = {
      name = lib.mkOption {
        type = types.str;
      };
      hostPath = lib.mkOption {
        type = types.str;
      };
      apiKeyFile = lib.mkOption {
        type = types.path;
      };
      maxDays = lib.mkOption {
        type = types.nullOr types.int;
        default = null;
      };
      maxSize = lib.mkOption {
        type = types.nullOr types.str;
        default = null;
      };
      includeTags = lib.mkOption {
        type = types.listOf types.str;
        default = [ ];
      };
      excludeTags = lib.mkOption {
        type = types.listOf types.str;
        default = [ ];
      };
    };
  };
  yaml = pkgs.formats.yaml { };
  filterOutNulls = x: lib.attrsets.filterAttrsRecursive (n: v: v != null) x;
  cleanedConfig = cfg.settings // {
    sonarr = map filterOutNulls cfg.settings.sonarr;
    radarr = map filterOutNulls cfg.settings.radarr;
  };
  configFile = yaml.generate "config.yaml" cleanedConfig;
in
{
  options.services.cleanarr = {
    enable = lib.mkEnableOption "Enable cleanarr";
    settings = lib.mkOption {
      default = { };
      type = types.submodule {
        options = {
          interval = lib.mkOption {
            type = types.int;
            default = 1;
          };
          dryRun = lib.mkOption {
            type = types.bool;
            default = false;
          };
          radarr = lib.mkOption {
            type = types.listOf servarrSettingType;
            default = [ ];
          };
          sonarr = lib.mkOption {
            type = types.listOf servarrSettingType;
            default = [ ];
          };
        };
      };
    };
  };

  config = lib.mkIf cfg.enable {
    systemd.services.cleanarr = {
      description = "Cleanarr";
      after = [ "network.target" ];
      wantedBy = [ "default.target" ];
      environment = {
        CLEANARR_CONFIG = configFile;
      };
      serviceConfig = {
        DynamicUser = true;
        ExecStart = lib.getExe cleanarr;
      };
    };
  };
}
