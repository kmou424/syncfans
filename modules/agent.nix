{
  config,
  lib,
  pkgs,
  ...
}:

with lib;

let
  cfg = config.services.syncfans-agent;

  # Generate TOML configuration from Nix options
  generateConfig = ''
    [config]
    server = "${cfg.server}"
    secret = "${cfg.secret}"
    fan = "${cfg.fan}"
    critical_temp_range = [${concatStringsSep ", " (map toString cfg.criticalTempRange)}]
    critical_margin = ${toString cfg.criticalMargin}
    override_curve = ${if cfg.overrideCurve then "true" else "false"}
    curve_type = "${cfg.curve.type}"
    curve_factor = ${toString cfg.curve.factor}
    dead_zone_ratio = ${toString cfg.curve.deadZoneRatio}
  ''
  + concatStringsSep "\n" (
    mapAttrsToList (sysinfoName: sysinfo: ''
      [sysinfo.${sysinfoName}]
      method = "${sysinfo.method}"
      query = "${sysinfo.query}"
      type = "${sysinfo.type}"
    '') cfg.sysinfo
  );

  configFile =
    if cfg.configFile != null then cfg.configFile else pkgs.writeText "agent.toml" generateConfig;
in
{
  options.services.syncfans-agent = {
    enable = mkEnableOption "SyncFans agent";

    package = mkOption {
      type = types.package;
      default =
        pkgs.syncfans-agent
          or (throw "syncfans-agent package not found. Please specify services.syncfans-agent.package or add syncfans to your flake inputs.");
      defaultText = "pkgs.syncfans-agent";
      description = "The SyncFans agent package to use.";
    };

    server = mkOption {
      type = types.str;
      description = "Server URL for reporting.";
      example = "http://127.0.0.1:16380/report";
    };

    secret = mkOption {
      type = types.str;
      description = "Secret key for authentication.";
    };

    fan = mkOption {
      type = types.str;
      description = "Fan name to control.";
      example = "fan1";
    };

    criticalTempRange = mkOption {
      type = types.listOf types.float;
      default = [
        40.0
        75.0
      ];
      description = "Critical temperature range in Celsius [min, max].";
    };

    criticalMargin = mkOption {
      type = types.float;
      default = 3.0;
      description = "Critical temperature margin in Celsius.";
    };

    overrideCurve = mkOption {
      type = types.bool;
      default = false;
      description = "Whether to override server curve parameters.";
    };

    curve = {
      type = mkOption {
        type = types.enum [
          "linear"
          "s-curve"
          "exponential"
          "aggressive"
        ];
        default = "s-curve";
        description = "Fan curve type.";
      };

      factor = mkOption {
        type = types.float;
        default = 1.0;
        description = "Curve factor.";
      };

      deadZoneRatio = mkOption {
        type = types.float;
        default = 0.1;
        description = "Dead zone ratio.";
      };
    };

    sysinfo = mkOption {
      type = types.attrsOf (
        types.submodule {
          options = {
            method = mkOption {
              type = types.enum [
                "shell"
                "file"
              ];
              description = "Query method.";
            };

            query = mkOption {
              type = types.str;
              description = "Shell command or file path.";
              example = "nvidia-smi --query-gpu=temperature.gpu --format=csv,noheader";
            };

            type = mkOption {
              type = types.enum [
                "float"
                "int"
                "string"
              ];
              description = "Value type.";
            };
          };
        }
      );
      default = {
        temperature = {
          method = "shell";
          query = "nvidia-smi --query-gpu=temperature.gpu --format=csv,noheader";
          type = "float";
        };
        usage = {
          method = "shell";
          query = "nvidia-smi --query-gpu=utilization.gpu --format=csv,noheader";
          type = "float";
        };
      };
      description = "System information queries.";
    };

    configFile = mkOption {
      type = types.nullOr types.path;
      default = null;
      description = "Path to external TOML config file (overrides generated config).";
    };

    environment = mkOption {
      type = types.attrsOf types.str;
      default = {
        SYNCFANS_DEBUG = "false";
      };
      description = "Environment variables for the service.";
    };

    requiresNvidia = mkOption {
      type = types.bool;
      default = true;
      description = "Whether to require nvidia-smi (adds it to PATH).";
    };
  };

  config = mkIf cfg.enable {
    systemd.services.syncfans-agent = {
      description = "SyncFans Agent";
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ];

      serviceConfig = {
        Type = "simple";
        ExecStart = "${cfg.package}/bin/agent -config ${configFile}";
        Restart = "always";
        RestartSec = "10s";
        # If agent disconnects, restart after 10 seconds to ensure server cleans up
        ExecStop = "${pkgs.coreutils}/bin/sleep 10";
        # Set PATH environment variable (Path option is deprecated in systemd)
        # Include bash for sh command execution, nvidia-smi, and system paths
        Environment =
          let
            basePath = makeBinPath [
              pkgs.bash
              pkgs.coreutils
              pkgs.findutils
              pkgs.gnugrep
              pkgs.gnused
              pkgs.systemd
            ];
            systemPath = "/run/current-system/sw/bin";
            nvidiaPath = optionalString (
              cfg.requiresNvidia && config.hardware.nvidia.enabled
            ) "${config.hardware.nvidia.package.bin}/bin:";
            fullPath = "${nvidiaPath}${basePath}:${systemPath}";
            envVars = cfg.environment // {
              PATH = fullPath;
            };
          in
          mapAttrsToList (name: value: "${name}=${value}") envVars;
      };
    };
  };
}
