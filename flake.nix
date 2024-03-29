{
  description = "ODALA";

  inputs.nixpkgs.url = "nixpkgs/release-22.05";

  outputs = { self, nixpkgs }: with nixpkgs.legacyPackages.x86_64-linux; {
    devShell.x86_64-linux = mkShell {
      buildInputs = [ go_1_18 jq golangci-lint ];
    };
  };
}
