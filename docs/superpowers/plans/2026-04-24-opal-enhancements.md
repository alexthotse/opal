# Opal Phase 1: TUI & Core Reasoning Enhancements Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the three distinct TUI themes, integrate Harmonica for spring animations in the Go frontend, and implement the ULTRATHINK and ULTRAPLAN feature flags in the Gleam backend.

**Architecture:** The Go frontend will use a new `theme.go` file to handle styling and `animator.go` using `charmbracelet/harmonica` for physics-based rendering. The Gleam backend will expand its JSON-RPC handler to support `ultrathink.start` and `ultraplan.start`.

**Tech Stack:** Go 1.25, Bubble Tea, Lipgloss, Harmonica, Gleam 1.9.0, Gleeunit.

---

### Task 1: Add Harmonica Dependency and Setup Animation Base

**Files:**
- Modify: `opal_frontend/go.mod`
- Create: `opal_frontend/animator.go`

- [ ] **Step 1: Get the harmonica dependency**

```bash
cd opal_frontend
go get github.com/charmbracelet/harmonica
```

- [ ] **Step 2: Create the animator file with a basic spring model**

```go
package main

import (
	"github.com/charmbracelet/harmonica"
)

type Animator struct {
	spring   harmonica.Spring
	position float64
	velocity float64
	target   float64
}

func NewAnimator() *Animator {
	return &Animator{
		spring:   harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.5),
		position: 0,
		velocity: 0,
		target:   0,
	}
}

func (a *Animator) Update() {
	a.position, a.velocity = a.spring.Update(a.position, a.velocity, a.target)
}

func (a *Animator) SetTarget(t float64) {
	a.target = t
}
```

- [ ] **Step 3: Commit**

```bash
git add opal_frontend/go.mod opal_frontend/go.sum opal_frontend/animator.go
git commit -m "feat: add harmonica for physics-based TUI animations"
```

### Task 2: Implement Theme Manager (PI.dev, free-code, crush)

**Files:**
- Create: `opal_frontend/theme.go`

- [ ] **Step 1: Write the theme definitions**

```go
package main

import "github.com/charmbracelet/lipgloss"

type ThemeType string

const (
	ThemePI       ThemeType = "pi.dev"
	ThemeFreeCode ThemeType = "free-code"
	ThemeCrush    ThemeType = "crush"
)

type Theme struct {
	Primary   lipgloss.Style
	Secondary lipgloss.Style
	Text      lipgloss.Style
}

func GetTheme(t ThemeType) Theme {
	switch t {
	case ThemeFreeCode:
		return Theme{
			Primary:   lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("255")).Bold(true),
			Secondary: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
			Text:      lipgloss.NewStyle().Foreground(lipgloss.Color("255")),
		}
	case ThemeCrush:
		return Theme{
			Primary:   lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true),
			Secondary: lipgloss.NewStyle().Foreground(lipgloss.Color("99")),
			Text:      lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
		}
	default: // ThemePI
		return Theme{
			Primary:   lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true),
			Secondary: lipgloss.NewStyle().Foreground(lipgloss.Color("43")),
			Text:      lipgloss.NewStyle().Foreground(lipgloss.Color("250")),
		}
	}
}
```

- [ ] **Step 2: Commit**

```bash
git add opal_frontend/theme.go
git commit -m "feat: add theme manager for PI.dev, free-code, and crush styles"
```

### Task 3: Integrate Theme and Animator into Main UI

**Files:**
- Modify: `opal_frontend/main.go`

- [ ] **Step 1: Update the model struct and View**

```bash
sed -i 's/type model struct {/type model struct {\n\ttheme      Theme\n\tanimator   *Animator\n/' opal_frontend/main.go
```

- [ ] **Step 2: Initialize theme and animator in `initialModel`**

```bash
sed -i 's/return model{/return model{\n\t\ttheme:      GetTheme(ThemePI),\n\t\tanimator:   NewAnimator(),\n/' opal_frontend/main.go
```

- [ ] **Step 3: Update `View` to use the theme**

```bash
sed -i 's/s := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render("✦ Opal Agent (Go + Gleam)")/s := m.theme.Primary.Render("✦ Opal Agent (" + string(ThemePI) + " theme)")/' opal_frontend/main.go
```

- [ ] **Step 4: Commit**

```bash
git add opal_frontend/main.go
git commit -m "feat: integrate theme and animator into main UI loop"
```

### Task 4: Add ULTRATHINK and ULTRAPLAN to Gleam Backend

**Files:**
- Modify: `opal_backend/src/opal_backend.gleam`
- Modify: `opal_backend/test/opal_backend_test.gleam`

- [ ] **Step 1: Write the failing test for new RPC methods**

```gleam
// Append to opal_backend/test/opal_backend_test.gleam

pub fn ultrathink_test() {
  opal_backend.execute_method("ultrathink.start")
  |> should.equal("ultrathink_mode_activated")
}

pub fn ultraplan_test() {
  opal_backend.execute_method("ultraplan.start")
  |> should.equal("ultraplan_mode_activated")
}
```

- [ ] **Step 2: Run the test to verify it fails**

```bash
cd opal_backend && gleam test
```
Expected: FAIL (unknown_method)

- [ ] **Step 3: Export `execute_method` and add the new methods**

```bash
sed -i 's/fn execute_method/pub fn execute_method/' opal_backend/src/opal_backend.gleam
```

```gleam
// Replace the execute_method function in opal_backend/src/opal_backend.gleam
pub fn execute_method(method: String) -> String {
  case method {
    "ping" -> "pong"
    "agent.start" -> "agent_started_successfully"
    "agent.status" -> "idle"
    "ultrathink.start" -> "ultrathink_mode_activated"
    "ultraplan.start" -> "ultraplan_mode_activated"
    _ -> "unknown_method"
  }
}
```

- [ ] **Step 4: Run the test to verify it passes**

```bash
cd opal_backend && gleam test
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add opal_backend/src/opal_backend.gleam opal_backend/test/opal_backend_test.gleam
git commit -m "feat: add ultrathink and ultraplan RPC methods to backend"
```
