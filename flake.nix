{
  description = "stonkler - a Go CLI for pluggable financial data backends";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs, ... }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      forAllSystems =
        f:
        nixpkgs.lib.genAttrs systems (
          system:
          f (
            import nixpkgs {
              inherit system;
            }
          )
        );
    in
    {
      packages = forAllSystems (pkgs: {
        stonk = pkgs.buildGoModule {
          pname = "stonk";
          version = "0.1.0";
          src = pkgs.lib.cleanSource ./.;
          vendorHash = "sha256-0Qxw+MUYVgzgWB8vi3HBYtVXSq/btfh4ZfV/m1chNrA=";
          subPackages = [ "cmd/stonk" ];
        };

        default = pkgs.buildGoModule {
          pname = "stonk";
          version = "0.1.0";
          src = pkgs.lib.cleanSource ./.;
          vendorHash = "sha256-0Qxw+MUYVgzgWB8vi3HBYtVXSq/btfh4ZfV/m1chNrA=";
          subPackages = [ "cmd/stonk" ];
        };
      });

      apps = forAllSystems (pkgs: {
        stonk = {
          type = "app";
          program = "${self.packages.${pkgs.stdenv.hostPlatform.system}.stonk}/bin/stonk";
        };

        default = {
          type = "app";
          program = "${self.packages.${pkgs.stdenv.hostPlatform.system}.default}/bin/stonk";
        };
      });

      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            cacert
            curl
            delve
            git
            go
            go-tools
            golangci-lint
            gopls
            gotools
            jq
            just
            nixfmt
          ];

          env = {
            CGO_ENABLED = "0";
            GOFLAGS = "-mod=mod";
          };

          shellHook = ''
            echo "stonkler dev shell"
            echo "Go: $(go version)"
            echo "Set FMP_API_KEY to use the Financial Modeling Prep backend."
          '';
        };
      });

      formatter = forAllSystems (
        pkgs:
        pkgs.writeShellApplication {
          name = "stonkler-format";
          runtimeInputs = [ pkgs.nixfmt ];
          text = ''
            if [ "$#" -eq 0 ]; then
              exec nixfmt flake.nix
            fi

            exec nixfmt "$@"
          '';
        }
      );
    };
}
