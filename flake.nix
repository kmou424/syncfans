{
  description = "SyncFans - Fan speed synchronization for GPU passthrough in PVE";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };

        syncfans-server = pkgs.buildGoModule {
          pname = "syncfans-server";
          version = "0.1.0";
          src = ./.;

          vendorHash = "sha256-8bg6caSu2udAEejh3PbjiUjemAA3qGLj4AgVxej/LO4=";

          subPackages = [ "cmd/server" ];

          meta = with pkgs.lib; {
            description = "SyncFans server - Fan speed control server";
            homepage = "https://github.com/kmou424/syncfans";
            license = licenses.mit;
            maintainers = [ ];
          };
        };

        syncfans-agent = pkgs.buildGoModule {
          pname = "syncfans-agent";
          version = "0.1.0";
          src = ./.;

          vendorHash = "sha256-8bg6caSu2udAEejh3PbjiUjemAA3qGLj4AgVxej/LO4=";

          subPackages = [ "cmd/agent" ];

          meta = with pkgs.lib; {
            description = "SyncFans agent - GPU temperature monitoring agent";
            homepage = "https://github.com/kmou424/syncfans";
            license = licenses.mit;
            maintainers = [ ];
          };
        };
      in
      {
        packages = {
          default = syncfans-server;
          server = syncfans-server;
          agent = syncfans-agent;
        };

        nixosModules = {
          server = import ./modules/server.nix;
          agent = import ./modules/agent.nix;
          default = import ./modules/server.nix;
        };
      }
    )
    // {
      nixosModules = {
        server = import ./modules/server.nix;
        agent = import ./modules/agent.nix;
        default = import ./modules/server.nix;
      };
    };
}
