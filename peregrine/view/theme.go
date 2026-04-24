package view

import (
	"github.com/charmbracelet/lipgloss"
)

type ThemeType int

const (
	ThemePIDev ThemeType = iota
	ThemeFreeCode
	ThemeCrush
)

type Theme struct {
	Primary lipgloss.Style
	Text    lipgloss.Style
	Border  lipgloss.Style
}

func GetTheme(t ThemeType) Theme {
	switch t {
	case ThemePIDev:
		return Theme{
			Primary: lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true),
			Text:    lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
			Border:  lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		}
	case ThemeFreeCode:
		return Theme{
			Primary: lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true),
			Text:    lipgloss.NewStyle().Foreground(lipgloss.Color("255")),
			Border:  lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
		}
	case ThemeCrush:
		return Theme{
			Primary: lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Italic(true),
			Text:    lipgloss.NewStyle().Foreground(lipgloss.Color("229")),
			Border:  lipgloss.NewStyle().Foreground(lipgloss.Color("141")),
		}
	default:
		return GetTheme(ThemePIDev)
	}
}
