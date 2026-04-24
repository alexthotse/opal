package view

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
