{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs = inputs: let
    system = "x86_64-linux";
    pkgs = import inputs.nixpkgs {inherit system;};
  in {
    devShells.${system}.default = pkgs.mkShell {};
    packages.${system}.default = pkgs.buildGoModule {
      pname = "clare";
      version = "v0.1";
      src = ./.;
      vendorHash = "sha256-xGPzODAJOls8RyyYdoEbPqz63i4oex41QsF1HNpmWAc=";
    };
  };
}
