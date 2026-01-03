package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuLevel int

const (
	MainMenu MenuLevel = iota
	AnalysisMenu
	SettingsMenu
	HelpMenu
)

type EnhancedMenuModel struct {
	cursor      int
	choices     []string
	level       MenuLevel
	parentMenu  MenuLevel
	SelectedVal string
	Done        bool
	width       int
	height      int
	showHelp    bool
}

func NewMenuModel() EnhancedMenuModel {
	return EnhancedMenuModel{
		choices: getMainMenuChoices(),
		level:   MainMenu,
	}
}

func getMainMenuChoices() []string {
	return []string{
		"📊 Analyze Repository",
		"🔄 Compare Repositories",
		"⚙️  Settings",
		"ℹ️  Help",
		"📜 Recent Analyses",
		"❌ Exit",
	}
}

func getAnalysisMenuChoices() []string {
	return []string{
		"🔍 Quick Analysis",
		"📈 Detailed Analysis",
		"📊 Custom Report",
		"⬅️  Back to Main Menu",
	}
}

func getSettingsMenuChoices() []string {
	return []string{
		"🎨 Theme (Light/Dark/Custom)",
		"🔐 GitHub Token Settings",
		"📁 Export Preferences",
		"🔄 Reset to Defaults",
		"⬅️  Back to Main Menu",
	}
}

func getHelpMenuChoices() []string {
	return []string{
		"⌨️  Keyboard Shortcuts",
		"❓ FAQ",
		"📖 Getting Started",
		"🔗 Documentation",
		"⬅️  Back to Main Menu",
	}
}

func (m EnhancedMenuModel) Init() tea.Cmd { return nil }

func (m EnhancedMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.handleSelection()
		case "h", "?":
			m.showHelp = !m.showHelp
		case "esc", "b":
			if m.level != MainMenu {
				m.navigateBack()
			}
		}
	}
	return m, nil
}

func (m *EnhancedMenuModel) handleSelection() {
	switch m.level {
	case MainMenu:
		m.handleMainMenuSelection()
	case AnalysisMenu:
		m.handleAnalysisMenuSelection()
	case SettingsMenu:
		m.handleSettingsMenuSelection()
	case HelpMenu:
		m.handleHelpMenuSelection()
	}
}

func (m *EnhancedMenuModel) handleMainMenuSelection() {
	switch m.cursor {
	case 0: // Analyze
		m.parentMenu = MainMenu
		m.level = AnalysisMenu
		m.choices = getAnalysisMenuChoices()
		m.cursor = 0
	case 1: // Compare
		m.SelectedVal = "compare"
		m.Done = true
	case 2: // Settings
		m.parentMenu = MainMenu
		m.level = SettingsMenu
		m.choices = getSettingsMenuChoices()
		m.cursor = 0
	case 3: // Help
		m.parentMenu = MainMenu
		m.level = HelpMenu
		m.choices = getHelpMenuChoices()
		m.cursor = 0
	case 4: // Recent
		m.SelectedVal = "recent"
		m.Done = true
	case 5: // Exit
		m.SelectedVal = "exit"
		m.Done = true
	}
}

func (m *EnhancedMenuModel) handleAnalysisMenuSelection() {
	switch m.cursor {
	case 0: // Quick Analysis
		m.SelectedVal = "quick_analyze"
		m.Done = true
	case 1: // Detailed Analysis
		m.SelectedVal = "detailed_analyze"
		m.Done = true
	case 2: // Custom Report
		m.SelectedVal = "custom_analyze"
		m.Done = true
	case 3: // Back
		m.navigateBack()
	}
}

func (m *EnhancedMenuModel) handleSettingsMenuSelection() {
	switch m.cursor {
	case 0: // Theme
		m.SelectedVal = "theme"
		m.Done = true
	case 1: // GitHub Token
		m.SelectedVal = "token"
		m.Done = true
	case 2: // Export Preferences
		m.SelectedVal = "export_prefs"
		m.Done = true
	case 3: // Reset
		m.SelectedVal = "reset"
		m.Done = true
	case 4: // Back
		m.navigateBack()
	}
}

func (m *EnhancedMenuModel) handleHelpMenuSelection() {
	switch m.cursor {
	case 0: // Shortcuts
		m.SelectedVal = "shortcuts"
		m.Done = true
	case 1: // FAQ
		m.SelectedVal = "faq"
		m.Done = true
	case 2: // Getting Started
		m.SelectedVal = "getting_started"
		m.Done = true
	case 3: // Documentation
		m.SelectedVal = "docs"
		m.Done = true
	case 4: // Back
		m.navigateBack()
	}
}

func (m *EnhancedMenuModel) navigateBack() {
	if m.level != MainMenu {
		m.level = m.parentMenu
		m.choices = getMainMenuChoices()
		m.cursor = 0
	}
}

func (m EnhancedMenuModel) View() string {
	logo := `
 ██████╗ ███████╗██████╗  ██████╗      ██╗     ██╗   ██╗ ███████╗  ███████╗██████╗ 
 ██╔══██╗██╔════╝██╔══██╗██╔═══██╗     ██║     ╚██╗ ██╔╝ ╚════██║  ██╔════╝██╔══██╗
 ██████╔╝█████╗  ██████╔╝██║   ██║█████╗██║      ╚████╔╝     ██╔╝  █████╗  ██████╔╝
 ██╔══██╗██╔══╝  ██╔═══╝ ██║   ██║╚════╝██║       ╚██╔╝     ██╔╝   ██╔══╝  ██╔══██╗
 ██║  ██║███████╗██║     ╚██████╔╝     ███████╗   ██║      ██╔╝     ███████╗██║  ██║   
 ╚═╝  ╚═╝╚══════╝╚═╝      ╚═════╝      ╚══════╝   ╚═╝     ███████╗ ╚══════╝╚═╝  ╚═╝     
`

	breadcrumb := m.getBreadcrumb()
	content := TitleStyle.Render(logo) + "\n\n" + SubtleStyle.Render(breadcrumb) + "\n\n"

	for i, choice := range m.choices {
		cursor := "  "
		style := NormalStyle

		if m.cursor == i {
			cursor = "▶ "
			style = SelectedStyle
		}

		content += fmt.Sprintf("%s%s\n", cursor, style.Render(choice))
	}

	// Footer with shortcuts
	footer := m.getFooter()
	content += "\n" + SubtleStyle.Render(footer)

	// Show contextual help if requested
	if m.showHelp {
		helpText := m.getContextualHelp()
		content += "\n\n" + HelpStyle.Render(helpText)
	}

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

func (m EnhancedMenuModel) getBreadcrumb() string {
	switch m.level {
	case MainMenu:
		return "📍 Main Menu"
	case AnalysisMenu:
		return "📍 Main Menu > Analysis Options"
	case SettingsMenu:
		return "📍 Main Menu > Settings"
	case HelpMenu:
		return "📍 Main Menu > Help"
	}
	return "📍 Main Menu"
}

func (m EnhancedMenuModel) getFooter() string {
	base := "↑ ↓ navigate • Enter select"
	if m.level != MainMenu {
		base += " • ESC back"
	}
	base += " • h help • q quit"
	return base
}

func (m EnhancedMenuModel) getContextualHelp() string {
	switch m.level {
	case MainMenu:
		return `
📝 NAVIGATION GUIDE:
  Analyze Repository - Quick or detailed analysis of GitHub repos
  Compare Repositories - Compare metrics between two repos
  Settings - Configure theme, GitHub token, and preferences
  Help - Keyboard shortcuts and documentation
  Recent Analyses - View your analysis history
  
🔑 QUICK KEYS: ? or h (help) • q (quit) • Enter (select)`
	case AnalysisMenu:
		return `
📊 ANALYSIS OPTIONS:
  Quick Analysis - Summary metrics and quick insights
  Detailed Analysis - Comprehensive analysis with all metrics
  Custom Report - Generate custom analysis report
  
⌨️  Back to main menu with ESC or select Back`
	case SettingsMenu:
		return `
⚙️  SETTINGS OPTIONS:
  Theme - Switch between light, dark, and custom themes
  GitHub Token - Configure GitHub API authentication
  Export Preferences - Set default export format and location
  Reset - Restore default settings
  
🔄 Settings are saved automatically`
	case HelpMenu:
		return `
❓ HELP RESOURCES:
  Keyboard Shortcuts - View all available keyboard shortcuts
  FAQ - Frequently asked questions
  Getting Started - Quick start guide
  Documentation - Full documentation and examples`
	}
	return ""
}
