package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HelpTopic int

const (
	HelpShortcuts HelpTopic = iota
	HelpFAQ
	HelpGettingStarted
	HelpDocumentation
)

type HelpModel struct {
	topic  HelpTopic
	cursor int
	width  int
	height int
	Done   bool
}

func NewHelpModel() HelpModel {
	return HelpModel{}
}

func (m HelpModel) Init() tea.Cmd { return nil }

func (m HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.Done = true
		case "1":
			m.topic = HelpShortcuts
		case "2":
			m.topic = HelpFAQ
		case "3":
			m.topic = HelpGettingStarted
		case "4":
			m.topic = HelpDocumentation
		}
	}
	return m, nil
}

func (m HelpModel) View() string {
	content := TitleStyle.Render("ℹ️  HELP CENTER") + "\n\n"

	// Topic selection
	topics := []string{
		"[1] Keyboard Shortcuts",
		"[2] FAQ - Frequently Asked Questions",
		"[3] Getting Started",
		"[4] Documentation",
	}

	for _, topic := range topics {
		content += NormalStyle.Render(topic) + "\n"
	}

	content += "\n" + SubtleStyle.Render("─────────────────────────────────────") + "\n\n"

	// Display selected topic
	helpContent := m.getHelpContent()
	content += helpContent + "\n\n"
	content += SubtleStyle.Render("Press ESC to return • q to quit")

	box := BoxStyle.Render(content)

	if m.width == 0 {
		return box
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Top,
		box,
	)
}

func (m HelpModel) getHelpContent() string {
	switch m.topic {
	case HelpShortcuts:
		return m.shortcutsHelp()
	case HelpFAQ:
		return m.faqHelp()
	case HelpGettingStarted:
		return m.gettingStartedHelp()
	case HelpDocumentation:
		return m.documentationHelp()
	}
	return m.defaultHelp()
}

func (m HelpModel) defaultHelp() string {
	return InfoStyle.Render("Select a help topic above to get started")
}

func (m HelpModel) shortcutsHelp() string {
	shortcuts := `
⌨️  KEYBOARD SHORTCUTS

Navigation:
  ↑/↓ or k/j       Navigate menu items
  Enter            Select menu item
  ESC/b            Go back / Back to menu
  q                Quit application

Analysis:
  1                Quick Analysis
  2                Detailed Analysis
  3                Custom Report

Menu:
  h/?              Show/hide help
  e                Export analysis
  j                Export as JSON
  m                Export as Markdown

Settings:
  s                Open Settings
  t                Change Theme

General:
  Ctrl+C           Quit application
  Ctrl+L           Clear screen
`
	return InfoStyle.Render(shortcuts)
}

func (m HelpModel) faqHelp() string {
	faq := `
❓ FREQUENTLY ASKED QUESTIONS

Q: How do I analyze a repository?
A: From the main menu, select "Analyze Repository" and choose your analysis type
   (Quick, Detailed, or Custom).

Q: What is the Bus Factor?
A: The Bus Factor measures repository risk by analyzing if the project relies too
   heavily on one or few contributors.

Q: How often is my GitHub token needed?
A: Your GitHub token is stored locally and used for API requests. You can configure
   it in Settings > GitHub Token Settings.

Q: Can I export my analysis?
A: Yes! After analysis, press 'e' to open export options. You can export as JSON
   or Markdown format.

Q: What are the health metrics?
A: Health score evaluates commit activity, issue handling, and contributor diversity.
   Higher scores indicate healthier repositories.

Q: Is my data private?
A: Yes! Repo-lyzer stores all data locally in ~/.repo-lyzer directory.
   Nothing is sent to external servers except GitHub API calls.
`
	return InfoStyle.Render(faq)
}

func (m HelpModel) gettingStartedHelp() string {
	guide := `
📖 GETTING STARTED

Step 1: Configure GitHub Token (Optional but recommended)
  1. Select "Settings" from main menu
  2. Choose "GitHub Token Settings"
  3. Enter your GitHub personal access token
  4. This increases API rate limits for faster analysis

Step 2: Analyze a Repository
  1. Select "Analyze Repository"
  2. Choose analysis type:
     - Quick: Summary metrics
     - Detailed: All available metrics
     - Custom: Choose specific metrics
  3. Enter repository name (format: owner/repo)
  4. Wait for analysis to complete

Step 3: View Results
  After analysis, view:
  - Health Score: Repository health status
  - Bus Factor: Dependency risk
  - Maturity Level: Project maturity assessment
  - Commit Activity: Recent commit trends

Step 4: Export Results
  - Press 'e' to export
  - Choose format (JSON or Markdown)
  - File is saved to default location in settings
`
	return InfoStyle.Render(guide)
}

func (m HelpModel) documentationHelp() string {
	docs := `
🔗 DOCUMENTATION

Repository Health Metrics:
  • Health Score: 0-100 scale measuring repository overall health
  • Commit Activity: Number of commits in time period
  • Contributors: Active contributors count
  • Languages: Primary programming languages used

Bus Factor Analysis:
  • Low Risk (7+): Healthy contributor distribution
  • Medium Risk (4-6): Some dependency concentration
  • High Risk (1-3): Critical dependency risk

Maturity Levels:
  • Experimental: New repository, minimal history
  • Development: Active development phase
  • Stable: Established, regular maintenance
  • Mature: Long-running, well-maintained project

Keyboard Modifiers:
  • Ctrl: Modifier key for system commands
  • Shift: Modifier for alternative actions (future)

For more info: https://github.com/agnivo988/Repo-lyzer
Report issues: https://github.com/agnivo988/Repo-lyzer/issues
`
	return InfoStyle.Render(docs)
}