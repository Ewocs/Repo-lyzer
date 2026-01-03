package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AnalysisHistory struct {
	Repository  string    `json:"repository"`
	Timestamp   time.Time `json:"timestamp"`
	HealthScore int       `json:"health_score"`
	Maturity    string    `json:"maturity"`
	FilePath    string    `json:"file_path"`
}

type HistoryManager struct {
	histories []AnalysisHistory
}

func getHistoryPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".repo-lyzer")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "history.json"), nil
}

func LoadHistory() ([]AnalysisHistory, error) {
	historyPath, err := getHistoryPath()
	if err != nil {
		return []AnalysisHistory{}, err
	}

	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		return []AnalysisHistory{}, nil
	}

	data, err := os.ReadFile(historyPath)
	if err != nil {
		return []AnalysisHistory{}, err
	}

	var histories []AnalysisHistory
	if err := json.Unmarshal(data, &histories); err != nil {
		return []AnalysisHistory{}, err
	}

	return histories, nil
}

func SaveHistory(histories []AnalysisHistory) error {
	historyPath, err := getHistoryPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(histories, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(historyPath, data, 0644)
}

func AddToHistory(result AnalysisResult, filePath string) error {
	histories, _ := LoadHistory()

	newEntry := AnalysisHistory{
		Repository:  result.Repo.FullName,
		Timestamp:   time.Now(),
		HealthScore: result.HealthScore,
		Maturity:    result.MaturityLevel,
		FilePath:    filePath,
	}

	// Keep only last 20 entries
	if len(histories) >= 20 {
		histories = histories[1:]
	}

	histories = append(histories, newEntry)
	return SaveHistory(histories)
}

type HistoryModel struct {
	histories []AnalysisHistory
	cursor    int
	width     int
	height    int
	selected  string
	Done      bool
}

func NewHistoryModel() HistoryModel {
	histories, _ := LoadHistory()
	return HistoryModel{
		histories: histories,
	}
}

func (m HistoryModel) Init() tea.Cmd { return nil }

func (m HistoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(m.histories)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor < len(m.histories) {
				m.selected = m.histories[m.cursor].Repository
				m.Done = true
			}
		case "esc":
			m.Done = true
		}
	}
	return m, nil
}

func (m HistoryModel) View() string {
	if len(m.histories) == 0 {
		return BoxStyle.Render("No recent analyses. Start by analyzing a repository!")
	}

	content := TitleStyle.Render("📜 RECENT ANALYSES") + "\n\n"

	for i := len(m.histories) - 1; i >= 0; i-- {
		h := m.histories[i]
		cursor := "  "
		style := NormalStyle

		if m.cursor == len(m.histories)-1-i {
			cursor = "▶ "
			style = SelectedStyle
		}

		timeStr := h.Timestamp.Format("2006-01-02 15:04")
		line := fmt.Sprintf("%s [%s] %s • Health: %d • Maturity: %s",
			cursor, timeStr, h.Repository, h.HealthScore, h.Maturity)

		content += style.Render(line) + "\n"
	}

	content += "\n" + SubtleStyle.Render("↑ ↓ navigate • Enter select • ESC back")

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