import sys, re, os
def main():
    if not os.path.exists("gleam.toml"):
        return
    with open("gleam.toml") as f:
        content = f.read()
    
    deps = []
    in_deps = False
    for line in content.splitlines():
        line = line.strip()
        if line.startswith("[dependencies]"):
            in_deps = True
            continue
        if line.startswith("[") and in_deps:
            in_deps = False
            continue
        if in_deps and "=" in line:
            dep_name = line.split("=")[0].strip()
            deps.append(f'"@hex_{dep_name}//:gleam"')
            
    deps_str = ",\n        ".join(deps)
    if deps:
        deps_str = "\n        " + deps_str + ",\n    "
        
    build_content = f"""load("@rules_gleam//gleam:defs.bzl", "gleam_library")

gleam_library(
    name = "gleam",
    srcs = glob([
        "**/*.gleam",
        "**/*.erl",
        "**/*.hrl",
        "**/*.app",
        "**/*.app.src",
        "gleam.toml",
    ], allow_empty = True),
    deps = [{deps_str}],
    visibility = ["//visibility:public"],
)
"""
    with open("BUILD.bazel", "w") as f:
        f.write(build_content)

if __name__ == "__main__":
    main()
