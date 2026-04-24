import sys
with open("/tmp/rules_gleam_patch/b/gleam_hex/repositories.bzl") as f:
    lines = f.readlines()

new_lines = []
skip = False
for line in lines:
    if 'cmd = [' in line:
        skip = True
        new_lines.append('''    ctx.file("generate_build.py", """import sys, os
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
        if in_deps and "=" in line and not line.startswith("#"):
            dep_name = line.split("=")[0].strip()
            deps.append(f'"@hex_{dep_name}//:gleam"')
            
    deps_str = ",\\\\n        ".join(deps)
    if deps:
        deps_str = "\\\\n        " + deps_str + ",\\\\n    "
    
    repo_name = os.path.basename(os.getcwd())
    strip_prefix = f"external/{repo_name}"
        
    build_content = f\\"\\"\\"load("@rules_gleam//gleam:defs.bzl", "gleam_library")

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
    strip_src_prefix = "{strip_prefix}",
    visibility = ["//visibility:public"],
)
\\"\\"\\"
    with open("BUILD.bazel", "w") as f:
        f.write(build_content)

if __name__ == "__main__":
    main()
""")
    result = env_execute(ctx, ["python3", "generate_build.py"])
''')
    elif skip and 'if result.return_code:' in line:
        skip = False
        new_lines.append('    if result.return_code:\n')
    elif skip:
        continue
    else:
        new_lines.append(line)

with open("/tmp/rules_gleam_patch/b/gleam_hex/repositories.bzl", "w") as f:
    f.writelines(new_lines)
