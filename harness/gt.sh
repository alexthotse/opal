#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
hq_dir="${GT_HQ_DIR:-"$root_dir/.gastown"}"
town_name="${GT_TOWN_NAME:-opal-hq}"
owner_email="${GT_OWNER_EMAIL:-ci@example.com}"
rig_name="${GT_RIG_NAME:-opal}"

need() {
  command -v "$1" >/dev/null 2>&1
}

ensure_path() {
  if need go; then
    local gopath
    gopath="$(go env GOPATH 2>/dev/null || true)"
    if [[ -n "${gopath:-}" ]]; then
      export PATH="$gopath/bin:$PATH"
    fi
  fi
}

install_tools() {
  ensure_path
  if ! need go; then
    echo "go missing"
    exit 1
  fi
  go install github.com/steveyegge/gastown/cmd/gt@latest
  go install github.com/steveyegge/beads/cmd/bd@latest
  ensure_path
  gt version >/dev/null
  bd version >/dev/null
}

init_hq() {
  ensure_path
  mkdir -p "$hq_dir"
  gt install "$hq_dir" --git --no-beads --owner "$owner_email" --name "$town_name" --force
  (
    cd "$hq_dir"
    bd init --non-interactive --role maintainer --prefix hq
    bd doctor
    if need dolt; then
      gt doctor --no-start
    fi
  )
}

add_rig() {
  ensure_path
  if [[ ! -d "$hq_dir" ]]; then
    init_hq
  fi
  (
    cd "$hq_dir"
    if gt rig list 2>/dev/null | grep -q "$rig_name"; then
      exit 0
    fi
    gt rig add "$rig_name" "$root_dir"
  )
}

create_bead() {
  ensure_path
  local title="${1:-}"
  if [[ -z "$title" ]]; then
    echo "title required"
    exit 1
  fi
  (
    cd "$hq_dir"
    bd q --title "$title"
  )
}

sling_bead() {
  ensure_path
  local bead_id="${1:-}"
  if [[ -z "$bead_id" ]]; then
    echo "bead id required"
    exit 1
  fi
  (
    cd "$hq_dir"
    gt sling "$bead_id" "$rig_name"
  )
}

doctor() {
  ensure_path
  (
    cd "$hq_dir"
    bd doctor || true
    gt doctor || true
  )
}

case "${1:-}" in
  install)
    install_tools
    ;;
  init)
    init_hq
    ;;
  rig-add)
    add_rig
    ;;
  bead)
    shift
    create_bead "$*"
    ;;
  sling)
    shift
    sling_bead "${1:-}"
    ;;
  doctor)
    doctor
    ;;
  *)
    echo "usage: $0 install|init|rig-add|bead <title>|sling <bead_id>|doctor"
    exit 2
    ;;
esac
