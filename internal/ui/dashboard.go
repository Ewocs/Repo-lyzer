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
		return m, tea.Tick(3*time.Second, func(time.Time) tea.Msg {
			return "clear_status"
		})

	case string:
		if msg == "clear_status" {
			m.statusMsg = ""
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.BackToMenu = true

		case "e":
			m.showExport = !m.showExport

		case "j":
			if m.showExport {
				return m, func() tea.Msg {
					// Stub export — replace with real ExportJSON later
					return exportMsg{nil, "Exported to analysis.json"}
				}
			}

		case "f":
			return m, func() tea.Msg { return "switch_to_tree" }
		}
	}

	return m, nil
}

func (m DashboardModel) View() string {
	if m.data.Repo == nil {
		return "No data loaded"
	}

	header := TitleStyle.Render(
		fmt.Sprintf("📊 Analysis for %s", m.data.Repo.FullName),
	)

	metrics := fmt.Sprintf(
		"Health Score: %d\nBus Factor: %d (%s)\nMaturity: %s (%d)",
		m.data.HealthScore,
		m.data.BusFactor,
		m.data.BusRisk,
		m.data.MaturityLevel,
		m.data.MaturityScore,
	)
	metricsBox := BoxStyle.Render(metrics)

	activity := analyzer.CommitsPerDay(m.data.Commits)
	chart := RenderCommitActivity(activity, 10)
	chartBox := BoxStyle.Render(chart)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		lipgloss.JoinHorizontal(lipgloss.Top, metricsBox, chartBox),
	)

	if m.showExport {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			BoxStyle.Render("📥 Export:\n[J] JSON"),
		)
	}

	if m.statusMsg != "" {
		content += "\n" + SubtleStyle.Render(m.statusMsg)
	}

	footer := SubtleStyle.Render("e: export • f: file tree • q: back")
	content += "\n" + footer

	if m.width == 0 {
		return content
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
