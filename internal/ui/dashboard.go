package ui

import (
	"fmt"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DashboardModel struct {
	data       AnalysisResult
	BackToMenu bool
	width      int
	height     int
	showExport bool
	statusMsg  string
}

func NewDashboardModel() DashboardModel {
	return DashboardModel{}
}

func (m DashboardModel) Init() tea.Cmd { return nil }

func (m *DashboardModel) SetData(data AnalysisResult) {
	m.data = data
}

type exportMsg struct {
	err error
	msg string
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case exportMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Export failed: %v", msg.err)
		} else {
			m.statusMsg = msg.msg
		}
		// Clear status after 3 seconds
		return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg { return "clear_status" })

	case string:
		if msg == "clear_status" {
			m.statusMsg = ""
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			if m.showExport {
				m.showExport = false
			} else {
				m.BackToMenu = true
			}
		case "e":
			m.showExport = !m.showExport
		case "j":
			if m.showExport {
				return m, func() tea.Msg {
					err := ExportJSON(m.data, "analysis.json")
					return exportMsg{err, "Exported to analysis.json"}
				}
			}
		} else {
			switch msg.String() {
			case "esc", "q":
				m.BackToMenu = true
			case "e":
				m.exportMenuVisible = !m.exportMenuVisible
				m.exportCursor = 0
			case "f":
				// Switch to tree view - handled in app.go
				return m, func() tea.Msg { return "switch_to_tree" }()
			}
		}
	}
	return m, nil
}

func (m DashboardModel) View() string {
	if m.data.Repo == nil {
		return "No data"
	}

	// Header
	header := TitleStyle.Render(fmt.Sprintf("Analysis for %s", m.data.Repo.FullName))

	// Metrics Column
	metrics := fmt.Sprintf(
		"Health Score: %d\nBus Factor: %d (%s)\nMaturity: %s (%d)",
		m.data.HealthScore,
		m.data.BusFactor, m.data.BusRisk,
		m.data.MaturityLevel, m.data.MaturityScore,
	)
	metricsBox := BoxStyle.Render(metrics)

	// Charts
	activityData := analyzer.CommitsPerDay(m.data.Commits)
	chart := RenderCommitActivity(activityData, 10) // Show last 10 days
	chartBox := BoxStyle.Render(chart)

	// Languages
	langs := "Languages:\n"
	for l, b := range m.data.Languages {
		langs += fmt.Sprintf("• %s (%d bytes)\n", l, b)
	}
	langBox := BoxStyle.Render(langs)

	// Top Contributors
	contribs := "Top Contributors:\n"
	limit := 5
	if len(m.data.Contributors) < 5 {
		limit = len(m.data.Contributors)
	}
	for i := 0; i < limit; i++ {
		c := m.data.Contributors[i]
		contribs += fmt.Sprintf("• %s (%d)\n", c.Login, c.Contributions)
	}
	contribBox := BoxStyle.Render(contribs)

	// Layout
	row1 := lipgloss.JoinHorizontal(lipgloss.Top, metricsBox, chartBox)
	row2 := lipgloss.JoinHorizontal(lipgloss.Top, langBox, contribBox)
	content := lipgloss.JoinVertical(lipgloss.Left, header, SubtleStyle.Render(repoInfo), row1, row2)

	if m.exportMenuVisible {
		exportMenu := "📥 EXPORT OPTIONS:\n\n"
		for i, opt := range m.exportOptions {
			cursor := "  "
			style := NormalStyle
			if m.exportCursor == i {
				cursor = "▶ "
				style = SelectedStyle
			}
			exportMenu += fmt.Sprintf("%s%s\n", cursor, style.Render(opt))
		}
		treeContent += fmt.Sprintf("%s %s\n", icon, m.data.FileTree[i].Path)
	}
	if len(m.data.FileTree) > limit {
		treeContent += fmt.Sprintf("... and %d more", len(m.data.FileTree)-limit)
	}
	treeBox := BoxStyle.Render(treeContent)

	// Layout
	row1 := lipgloss.JoinHorizontal(lipgloss.Top, metricsBox, chartBox)
	content := lipgloss.JoinVertical(lipgloss.Left, header, row1, treeBox)

	if m.showExport {
		exportMenu := BoxStyle.Render("Export Options:\n[J] JSON\n[M] Markdown")
		content = lipgloss.JoinVertical(lipgloss.Left, content, exportMenu)
	}

	footer := SubtleStyle.Render("e: export • f: file tree")
	if !m.exportMenuVisible {
		footer += SubtleStyle.Render(" • q: back")
	} else {
		footer += SubtleStyle.Render(" • ↑ ↓ select • Enter confirm • ESC close")
	}

	content += "\n" + SubtleStyle.Render("e: export • q: back")

	if m.width == 0 {
		return content
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}
