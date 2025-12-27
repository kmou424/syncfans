{
  config,
  lib,
  pkgs,
  ...
}:

with lib;

let
  cfg = config.services.syncfans-server;

  # Generate TOML configuration from Nix options
  generateConfig = ''
    [config]
    listen = "${cfg.listen}"
    secret = "${cfg.secret}"
    interval = ${toString cfg.interval}
    smoothing_factor = ${toString cfg.smoothingFactor}

    [default]
    curve_type = "${cfg.defaultCurve.type}"
    curve_factor = ${toString cfg.defaultCurve.factor}
    dead_zone_ratio = ${toString cfg.defaultCurve.deadZoneRatio}
  ''
  + concatStringsSep "\n" (
    mapAttrsToList (fanName: fan: ''
      [sysfans.${fanName}]
      path = "${fan.path}"
      max_speed = ${toString fan.maxSpeed}
      min_speed = ${toString fan.minSpeed}
    '') cfg.fans
  );

  configFile =
    if cfg.configFile != null then cfg.configFile else pkgs.writeText "server.toml" generateConfig;
in
{
  options.services.syncfans-server = {
    enable = mkEnableOption "SyncFans server";

    package = mkOption {
      type = types.package;
      default =
        pkgs.syncfans-server
          or (throw "syncfans-server package not found. Please specify services.syncfans-server.package or add syncfans to your flake inputs.");
      defaultText = "pkgs.syncfans-server";
      description = "The SyncFans server package to use.";
    };

    listen = mkOption {
      type = types.str;
      default = "127.0.0.1:16380";
      description = "Listen address for the server.";
    };

    secret = mkOption {
      type = types.str;
      description = "Secret key for authentication.";
    };

    interval = mkOption {
      type = types.int;
      default = 1000;
      description = "Report processing interval in milliseconds.";
    };

    smoothingFactor = mkOption {
      type = types.float;
      default = 0.2;
      description = "Smoothing factor for fan control (smaller value = smoother).";
    };

    defaultCurve = {
      type = mkOption {
        type = types.enum [
          "linear"
          "s-curve"
          "exponential"
          "aggressive"
        ];
        default = "s-curve";
        description = "Default fan curve type.";
      };

      factor = mkOption {
        type = types.float;
        default = 1.8;
        description = "Curve factor (not used for linear and s-curve).";
      };

      deadZoneRatio = mkOption {
        type = types.float;
        default = 0.1;
        description = "Dead zone ratio for temperature control.";
      };
    };

    fans = mkOption {
      type = types.attrsOf (
        types.submodule {
          options = {
            path = mkOption {
              type = types.str;
              description = "Sysfs path for the fan PWM control.";
              example = "/sys/class/hwmon/hwmon4/pwm1";
            };

            maxSpeed = mkOption {
              type = types.int;
              description = "Maximum fan speed (PWM value).";
              example = 255;
            };

            minSpeed = mkOption {
              type = types.int;
              description = "Minimum fan speed (PWM value, prevent stall).";
              example = 30;
            };
          };
        }
      );
      default = { };
      description = "Fan configurations.";
      example = {
        fan1 = {
          path = "/sys/class/hwmon/hwmon4/pwm1";
          maxSpeed = 255;
          minSpeed = 30;
        };
      };
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
  };

  config = mkIf cfg.enable {
    systemd.services.syncfans-server = {
      description = "SyncFans Server";
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ];

      serviceConfig = {
        Type = "simple";
        ExecStart = "${cfg.package}/bin/server -config ${configFile}";
        Restart = "always";
        RestartSec = "10s";
        Environment = mapAttrsToList (name: value: "${name}=${value}") cfg.environment;
      };
    };
  };
}
