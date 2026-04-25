load("@rules_gleam//gleam:provider.bzl", "GleamErlPackageInfo")

def _elixir_beam_impl(ctx):
    out = ctx.actions.declare_file("Elixir.Falcon.JidoAgent.beam")
    ctx.actions.run_shell(
        inputs = ctx.files.src,
        outputs = [out],
        command = "elixirc -o $(dirname %s) %s" % (out.path, ctx.files.src[0].path),
        use_default_shell_env = True,
    )
    return [
        DefaultInfo(files = depset([out])),
        GleamErlPackageInfo(
            module_names = ["Elixir.Falcon.JidoAgent"],
            erl_module = depset([]),
            beam_module = depset([out]),
            gleam_cache = depset([]),
            strip_src_prefix = "",
        )
    ]

elixir_beam = rule(
    implementation = _elixir_beam_impl,
    attrs = {
        "src": attr.label(allow_single_file = [".ex"]),
    }
)