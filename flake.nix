# SPDX-License-Identifier: Unlicense OR MIT
{
  description = "Gio build environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
      in {
        devShells = {
          default = with pkgs;
            mkShell (rec {
              packages = [ clang ]
                ++ (if stdenv.isLinux then [
                  vulkan-headers
                  libxkbcommon
                  xorg.libX11
                  xorg.libXcursor
                  xorg.libXfixes
                  libGL
                  pkg-config
                ] else
                  [ ]);
            } // (if stdenv.isLinux then {
              LD_LIBRARY_PATH = "${vulkan-loader}/lib";
            } else
              { }));
        };
      });
}
