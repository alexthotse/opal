package view

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/alexthotse/peregrine/domain"
)

// Render converts the pure state and position into a rendered string.
func Render(state domain.State, pos float64) string {
	theme := GetTheme(ThemeType(state.ThemeType))

	pad := int(pos)
	if pad < 0 {
		pad = 0
	}
	paddingStyle := lipgloss.NewStyle().PaddingLeft(pad)

	s := theme.Primary.Render("✦ Peregrine Agent") + "\n\n"
	for _, msg := range state.Messages {
		s += theme.Text.Render(msg) + "\n"
	}
	return paddingStyle.Render(s)
}
