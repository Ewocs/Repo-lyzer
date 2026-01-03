package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Theme string

const (
	DarkTheme   Theme = "dark"
	LightTheme  Theme = "light"
	CustomTheme Theme = "custom"
)

type AppSettings struct {
	Theme              Theme  `json:"theme"`
	GitHubToken        string `json:"github_token"`
	DefaultExportType  string `json:"default_export_type"`
	DefaultExportPath  string `json:"default_export_path"`
	ShowProgressBars   bool   `json:"show_progress_bars"`
	EnableNotifications bool  `json:"enable_notifications"`
}

var DefaultSettings = AppSettings{
	Theme:              DarkTheme,
	DefaultExportType:  "json",
	DefaultExportPath:  "./exports",
	ShowProgressBars:   true,
	EnableNotifications: true,
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".repo-lyzer")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}

func LoadSettings() (AppSettings, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return DefaultSettings, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		return DefaultSettings, SaveSettings(DefaultSettings)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultSettings, err
	}

	var settings AppSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return DefaultSettings, err
	}

	return settings, nil
}

func SaveSettings(settings AppSettings) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func ResetSettings() error {
	return SaveSettings(DefaultSettings)
}

func (s *AppSettings) GetThemeStyle() (lipglossStyle, lipglossStyle) {
	switch s.Theme {
	case LightTheme:
		return LightBoxStyle, lipgloss.NewStyle().Foreground(lipgloss.Color("#000000"))
	case CustomTheme:
		// Return custom theme (can be extended)
		return BoxStyle, NormalStyle
	default:
		return DarkBoxStyle, NormalStyle
	}
}

type SettingsModel struct {
	settings   AppSettings
	cursor     int
	choices    []string
	width      int
	height     int
	edited     bool
	statusMsg  string
}

func NewSettingsModel(settings AppSettings) SettingsModel {
	return SettingsModel{
		settings: settings,
		choices: []string{
			"🎨 Theme: " + string(settings.Theme),
			"🔐 GitHub Token: " + hideToken(settings.GitHubToken),
			"📁 Export Type: " + settings.DefaultExportType,
			"💾 Export Path: " + settings.DefaultExportPath,
			"⚙️  Show Progress Bars: " + boolToStr(settings.ShowProgressBars),
			"🔔 Notifications: " + boolToStr(settings.EnableNotifications),
			"🔄 Reset to Defaults",
			"✅ Save & Exit",
		},
	}
}

func hideToken(token string) string {
	if token == "" {
		return "[Not configured]"
	}
	if len(token) > 8 {
		return token[:4] + "****" + token[len(token)-4:]
	}
	return "****"
}

func boolToStr(b bool) string {
	if b {
		return "✓ Enabled"
	}
	return "✗ Disabled"
}

func (m SettingsModel) Init() tea.Cmd { return nil }

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.handleSettingSelection()
		}
	}
	return m, nil
}

func (m *SettingsModel) handleSettingSelection() {
	switch m.cursor {
	case 0: // Theme
		themes := []Theme{DarkTheme, LightTheme, CustomTheme}
		for i, t := range themes {
			if t == m.settings.Theme {
				nextTheme := themes[(i+1)%len(themes)]
				m.settings.Theme = nextTheme
				m.choices[0] = "🎨 Theme: " + string(nextTheme)
				m.edited = true
				break
			}
		}
	case 4: // Progress Bars
		m.settings.ShowProgressBars = !m.settings.ShowProgressBars
		m.choices[4] = "⚙️  Show Progress Bars: " + boolToStr(m.settings.ShowProgressBars)
		m.edited = true
	case 5: // Notifications
		m.settings.EnableNotifications = !m.settings.EnableNotifications
		m.choices[5] = "🔔 Notifications: " + boolToStr(m.settings.EnableNotifications)
		m.edited = true
	case 6: // Reset
		m.settings = DefaultSettings
		m.choices = NewSettingsModel(DefaultSettings).choices
		m.edited = true
		m.statusMsg = "Settings reset to defaults"
	case 7: // Save & Exit
		if m.edited {
			SaveSettings(m.settings)
			m.statusMsg = "Settings saved successfully!"
		}
	}
}

func (m SettingsModel) View() string {
	content := TitleStyle.Render("⚙️  SETTINGS") + "\n\n"

	for i, choice := range m.choices {
		cursor := "  "
		style := NormalStyle

		if m.cursor == i {
			cursor = "▶ "
			style = SelectedStyle
		}

		content += fmt.Sprintf("%s%s\n", cursor, style.Render(choice))
	}

	if m.statusMsg != "" {
		content += "\n" + SuccessStyle.Render(m.statusMsg)
	}

	content += "\n" + SubtleStyle.Render("↑ ↓ navigate • Enter select • q quit")

	box := BoxStyle.Render(content)

	if m.width == 0 {
		return box
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		box,
	)
}