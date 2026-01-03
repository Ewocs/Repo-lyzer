package ui

import (
	"encoding/json"
	"fmt"
	"os"
)

func ExportJSON(data AnalysisResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func ExportMarkdown(data AnalysisResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	md := fmt.Sprintf("# Analysis for %s\n\n", data.Repo.FullName)
	md += fmt.Sprintf("## Health Score: %d\n", data.HealthScore)
	md += fmt.Sprintf("## Bus Factor: %d (%s)\n", data.BusFactor, data.BusRisk)
	md += fmt.Sprintf("## Maturity: %s (%d)\n", data.MaturityLevel, data.MaturityScore)
	
	md += "\n## File Tree (Top 20)\n"
	limit := 20
	if len(data.FileTree) < limit {
		limit = len(data.FileTree)
	}
	for i := 0; i < limit; i++ {
		icon := "📄"
		if data.FileTree[i].Type == "tree" {
			icon = "📁"
		}
		md += fmt.Sprintf("- %s %s\n", icon, data.FileTree[i].Path)
	}

	_, err = file.WriteString(md)
	return err
}
