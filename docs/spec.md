# Opal System Enhancement Specification (C.A.V.E. System)

## 1. Context (C)
Opal is a terminal-based AI assistant system featuring a Gleam backend (`opal_backend`) and a Go Bubble Tea frontend (`opal_frontend`). The goal of this enhancement is to bring the Opal system to feature parity with the `paoloanzn/free-code` fork of Claude Code CLI. Additionally, the Go-based frontend must mimic the fluid, spring-based animations of the former React-based TUI using `charmbracelet/harmonica`, and support three distinct, selectable visual themes:
- **PI.dev inspired** (default)
- **free-code inspired** (unrestricted, telemetry-free raw style)
- **opencode/crush inspired** (Charm-based, sleek and colorful terminal style)

## 2. Architecture (A)

### Backend (Gleam - `opal_backend`)
The Gleam backend will handle the core logic, state management, and LLM integrations. We will introduce feature flags and new modules for each of the experimental `free-code` capabilities. The JSON-RPC communication protocol will be expanded to support new methods for:
- Agent planning and multi-agent remote reasoning (`ULTRAPLAN`, `ULTRATHINK`)
- Background task automation and triggers (`AGENT_TRIGGERS`, `AGENT_TRIGGERS_REMOTE`)
- Deep memory extraction and caching (`EXTRACT_MEMORIES`, `CACHED_MICROCOMPACT`, `TEAMMEM`)
- Ideation and verification (`VERIFICATION_AGENT`, `SHOT_STATS`)
- Environment and security bridges (`BRIDGE_MODE`, `BASH_CLASSIFIER`)

### Frontend (Go Bubble Tea - `opal_frontend`)
The Go frontend will manage user interaction, rendering, and theming. 
- **Animations:** We will integrate `github.com/charmbracelet/harmonica` to power spring-physics animations for view transitions, loading indicators, and message reveals, mimicking the React Ink animations from Claude Code.
- **Theming:** A `Theme` interface will be implemented using `lipgloss`, with three concrete implementations (PI, FreeCode, Crush). The theme will dictate color palettes, border styles, typography emphasis, and animation spring stiffness.

## 3. Variables & Values (V) - Feature Parity Matrix

The following parity matrix maps the features from the `free-code` codebase to the targeted `opal_backend` implementation.

| free-code Feature | Opal Implementation | Status | Description |
| --- | --- | --- | --- |
| **ULTRAPLAN** | `opal_backend/src/planning.gleam` | Planned | Remote multi-agent planning. |
| **ULTRATHINK** | `opal_backend/src/reasoning.gleam` | Planned | Deep thinking mode for complex reasoning tasks. |
| **VOICE_MODE** | `opal_frontend/voice.go` | Planned | Push-to-talk voice input and dictation support. |
| **TOKEN_BUDGET** | `opal_backend/src/budget.gleam` | Planned | Token budget tracking and usage warnings. |
| **HISTORY_PICKER** | `opal_frontend/history.go` | Planned | Interactive prompt history picker in the TUI. |
| **MESSAGE_ACTIONS** | `opal_frontend/actions.go` | Planned | Inline message action entrypoints in the UI. |
| **QUICK_SEARCH** | `opal_backend/src/search.gleam` | Planned | Prompt quick-search over codebase. |
| **SHOT_STATS** | `opal_backend/src/stats.gleam` | Planned | Shot-distribution and LLM confidence stats. |
| **VERIFICATION_AGENT**| `opal_backend/src/verification.gleam`| Planned | Agent dedicated to task and code validation. |
| **AGENT_TRIGGERS** | `opal_backend/src/triggers.gleam` | Planned | Local cron/trigger tools for background automation. |
| **EXTRACT_MEMORIES** | `opal_backend/src/memory.gleam` | Planned | Post-query automatic memory extraction. |
| **COMPACTION_REMINDERS**| `opal_backend/src/memory.gleam`| Planned | Smart reminders around context compaction. |
| **CACHED_MICROCOMPACT**| `opal_backend/src/cache.gleam` | Planned | Cached microcompact state through query flows. |
| **TEAMMEM** | `opal_backend/src/teammem.gleam` | Planned | Team-memory files and watcher hooks. |
| **BRIDGE_MODE** | `opal_backend/src/bridge.gleam` | Planned | IDE remote-control bridge (VS Code, JetBrains). |
| **BASH_CLASSIFIER** | `opal_backend/src/security.gleam`| Planned | Classifier-assisted bash permission decisions. |

## 4. Execution (E)

Execution will follow a test-driven, bite-sized task approach:
1. **Frontend Architecture Update:** Integrate `harmonica` for physics-based animations. Add the Theme Engine and implement the three themes (PI.dev, free-code, opencode/crush).
2. **Backend Feature Rollout:** Implement Gleam modules for the `free-code` parity matrix. Start with the core reasoning/planning agents (`ULTRATHINK`, `ULTRAPLAN`) and memory management (`EXTRACT_MEMORIES`, `TEAMMEM`).
3. **RPC Bridge Enhancement:** Expand the JSON-RPC interface between the Go frontend and Gleam backend to support the new feature states and data flows.
4. **Testing & Verification:** Ensure all backend logic is covered by `gleeunit` and frontend UI states are tested for layout integrity.
