def gleam_binary(name, src_dir, deps = []):
    """
    A simple wrapper rule to invoke `gleam build` using the Bazel sandbox.
    It expects the Gleam compiler and Erlang to be provided by the host (or Nix).
    """
    native.genrule(
        name = name,
        srcs = native.glob([src_dir + "/**"]),
        outs = [name + "_out"],
        cmd = """
            # Copy source to a temporary directory inside the sandbox
            cp -R $(location {src_dir}/gleam.toml) $$(dirname $(location {src_dir}/gleam.toml))/../tmp_build
            cd tmp_build
            
            # Download deps and build
            gleam deps download
            gleam build
            
            # Export the build artifact (Erlang app/ebin) to the output
            cp -R build/dev/erlang/falcon/ebin $@
        """.format(src_dir = src_dir),
        tools = [], # Relies on the PATH having `gleam` provided by Nix
        visibility = ["//visibility:public"],
    )
