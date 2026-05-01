# SPEC

## §G GOAL
terminal AI assistant system. TUI client (Go) talks to local backend (Gleam) over ConnectRPC.

## §C CONSTRAINTS
- langs: Go + Gleam. no Python.
- Bazel build/test primary. `bazelisk test //peregrine/...` ! pass.
- Go code pure (no CGO). audio capture via external `arecord` on Linux.
- Nix builds Alpine-based OCI image (dockerTools) for `peregrine`.
- no secrets in repo. keys via env only.

## §I INTERFACES
- cmd: `peregrine [--theme <pi|freecode|crush>] [--provider <anthropic|openai|...>]` → TUI.
- file: `debug.log` created in CWD by `peregrine` run.
- api: Falcon backend expected at `http://localhost:8080` (ConnectRPC). `peregrine` acts as client.
- env: `PEREGRINE_ARECORD_DEVICE` ? ALSA device passed to `arecord -D`.
- env: `OPENAI_API_KEY` / provider keys required for agent client (no default secret).

## §V INVARIANTS
V1: repo contains 0 `*.py` source files.
V2: `bazelisk test //falcon:falcon_test //peregrine/...` exit 0.
V3: `go test ./...` exit 0 in `peregrine/`.
V4: Linux voice capture uses `arecord` (`-f S16_LE -c 1 -r 16000`) not embedded audio lib.
V5: Nix container base image = Alpine, built via nix (no Dockerfile).

## §T TASKS
id|status|task|cites
T1|x|remove Python scripts + python-based build hooks|V1
T2|x|replace malgo recorder with `arecord` recorder|V4
T3|x|fix Bazel to pure Go toolchain + remove malgo dep|V2
T4|x|fix Nix Alpine OCI image build in CI|V5
T5|x|keep tests green (Bazel + Go)|V2,V3

## §B BUGS
id|date|cause|fix
