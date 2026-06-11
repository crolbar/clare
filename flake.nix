{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs = inputs: let
    systems = ["x86_64-linux" "aarch64-linux"];
    forEachSystem = inputs.nixpkgs.lib.genAttrs systems;
    pkgsFor = inputs.nixpkgs.legacyPackages;
  in {
    devShells = forEachSystem (system: let
      pkgs = pkgsFor.${system};
    in {
      default = pkgs.mkShell {};
    });

    packages = forEachSystem (system: let
      pkgs = pkgsFor.${system};
    in {
      default = pkgs.buildGoModule {
        pname = "clare";
        version = "v0.1";
        src = ./.;
        vendorHash = "sha256-xGPzODAJOls8RyyYdoEbPqz63i4oex41QsF1HNpmWAc=";
      };
    });
  };
}
