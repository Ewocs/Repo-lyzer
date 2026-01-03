package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	dateStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00E5FF"))
	countStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB000"))
)

func barColor(count, max int) lipgloss.Style {
	if max == 0 {
		return lipgloss.NewStyle()
	}

	ratio := float64(count) / float64(max)

	switch {
	case ratio >= 0.67:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6E00"))
	case ratio >= 0.34:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#E89149"))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#292C7B"))
	}
}

func RenderCommitActivity(data map[string]int, maxDays int) string {
	var sb strings.Builder
	sb.WriteString(TitleStyle.Render("📈 Commit Activity") + "\n")

	dates := make([]string, 0, len(data))
	for d := range data {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	if len(dates) > maxDays {
		dates = dates[len(dates)-maxDays:]
	}

	max := 0
	for _, d := range dates {
		if data[d] > max {
			max = data[d]
		}
	}

	for _, d := range dates {
		count := data[d]
		barLen := 0
		if max > 0 {
			barLen = int(float64(count) / float64(max) * 20)
		}

		bar := barColor(count, max).Render(strings.Repeat("█", barLen))
		sb.WriteString(fmt.Sprintf(
			"%s | %s %s\n",
			dateStyle.Render(d),
			bar,
			countStyle.Render(fmt.Sprintf("%d", count)),
		))
	}
	return sb.String()
}
