package ui

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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

	md += "## Repository Information\n"
	md += fmt.Sprintf("- **Repository**: %s\n", data.Repo.FullName)
	md += fmt.Sprintf("- **URL**: %s\n", data.Repo.HTMLURL)
	md += fmt.Sprintf("- **Description**: %s\n", data.Repo.Description)
	md += fmt.Sprintf("- **Stars**: %d\n", data.Repo.StargazersCount)
	md += fmt.Sprintf("- **Forks**: %d\n\n", data.Repo.ForksCount)

	md += "## Health Metrics\n"
	md += fmt.Sprintf("- **Health Score**: %d/100\n", data.HealthScore)
	md += fmt.Sprintf("- **Bus Factor**: %d (%s)\n", data.BusFactor, data.BusRisk)
	md += fmt.Sprintf("- **Maturity Level**: %s (Score: %d)\n", data.MaturityLevel, data.MaturityScore)
	md += fmt.Sprintf("- **Contributors**: %d\n", len(data.Contributors))
	md += fmt.Sprintf("- **Total Commits**: %d\n\n", len(data.Commits))

	if len(data.Languages) > 0 {
		md += "## Programming Languages\n"
		for lang, bytes := range data.Languages {
			md += fmt.Sprintf("- %s: %d bytes\n", lang, bytes)
		}
		md += "\n"
	}

	md += "## Analysis Summary\n"
	if data.HealthScore >= 70 {
		md += "✅ This repository has a **good health score** and appears well-maintained.\n"
	} else if data.HealthScore >= 50 {
		md += "⚠️  This repository has a **moderate health score** with room for improvement.\n"
	} else {
		md += "❌ This repository has a **low health score** and may need attention.\n"
	}

	if data.BusFactor <= 2 {
		md += "⚠️  **Warning**: This repository has a HIGH bus factor risk. Consider increasing contributor diversity.\n"
	}

	_, err = file.WriteString(md)
	return err
}

func ExportCSV(data AnalysisResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	headers := []string{
		"Metric",
		"Value",
	}
	writer.Write(headers)

	// Write data rows
	records := [][]string{
		{"Repository", data.Repo.FullName},
		{"URL", data.Repo.HTMLURL},
		{"Description", data.Repo.Description},
		{"Stars", strconv.Itoa(data.Repo.StargazersCount)},
		{"Forks", strconv.Itoa(data.Repo.ForksCount)},
		{"Health Score", strconv.Itoa(data.HealthScore)},
		{"Bus Factor", strconv.Itoa(data.BusFactor)},
		{"Bus Risk", data.BusRisk},
		{"Maturity Level", data.MaturityLevel},
		{"Maturity Score", strconv.Itoa(data.MaturityScore)},
		{"Contributors", strconv.Itoa(len(data.Contributors))},
		{"Total Commits", strconv.Itoa(len(data.Commits))},
	}

	for _, record := range records {
		writer.Write(record)
	}

	return nil
}

func ExportHTML(data AnalysisResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Repository Analysis Report</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 900px;
            margin: 0 auto;
            background-color: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1, h2 {
            color: #333;
        }
        .metrics {
            display: grid;
            grid-template-columns: repeat(3, 1fr);
            gap: 20px;
            margin: 20px 0;
        }
        .metric {
            padding: 15px;
            background-color: #f9f9f9;
            border-left: 4px solid #007bff;
            border-radius: 4px;
        }
        .metric-value {
            font-size: 24px;
            font-weight: bold;
            color: #007bff;
        }
        .metric-label {
            font-size: 14px;
            color: #666;
            margin-top: 5px;
        }
        .good { border-left-color: #28a745; }
        .warning { border-left-color: #ffc107; }
        .danger { border-left-color: #dc3545; }
        .good .metric-value { color: #28a745; }
        .warning .metric-value { color: #ffc107; }
        .danger .metric-value { color: #dc3545; }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }
        table th {
            background-color: #007bff;
            color: white;
            padding: 10px;
            text-align: left;
        }
        table td {
            padding: 10px;
            border-bottom: 1px solid #ddd;
        }
        table tr:hover {
            background-color: #f5f5f5;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📊 Repository Analysis Report</h1>
        <p><strong>Repository:</strong> <a href="%s">%s</a></p>
        <p><strong>Analysis Date:</strong> %s</p>

        <h2>Key Metrics</h2>
        <div class="metrics">
            <div class="metric %s">
                <div class="metric-value">%d</div>
                <div class="metric-label">Health Score</div>
            </div>
            <div class="metric %s">
                <div class="metric-value">%d</div>
                <div class="metric-label">Bus Factor</div>
            </div>
            <div class="metric">
                <div class="metric-value">%s</div>
                <div class="metric-label">Maturity</div>
            </div>
        </div>

        <h2>Repository Details</h2>
        <table>
            <tr>
                <th>Attribute</th>
                <th>Value</th>
            </tr>
            <tr>
                <td>Stars</td>
                <td>%d</td>
            </tr>
            <tr>
                <td>Forks</td>
                <td>%d</td>
            </tr>
            <tr>
                <td>Contributors</td>
                <td>%d</td>
            </tr>
            <tr>
                <td>Total Commits</td>
                <td>%d</td>
            </tr>
            <tr>
                <td>Bus Risk</td>
                <td>%s</td>
            </tr>
        </table>

        <h2>Summary</h2>
        <p>%s</p>
    </div>
</body>
</html>`

	healthClass := "good"
	if data.HealthScore < 70 {
		healthClass = "warning"
	}
	if data.HealthScore < 50 {
		healthClass = "danger"
	}

	busClass := "good"
	if data.BusFactor <= 2 {
		busClass = "danger"
	} else if data.BusFactor <= 4 {
		busClass = "warning"
	}

	summary := fmt.Sprintf("This repository has %d contributors and %d commits. ", len(data.Contributors), len(data.Commits))
	if data.BusFactor <= 2 {
		summary += "⚠️ High bus factor risk detected. "
	}
	if data.HealthScore >= 70 {
		summary += "✅ Good health score."
	} else {
		summary += "⚠️ Moderate to low health score."
	}

	finalHTML := fmt.Sprintf(html,
		data.Repo.HTMLURL, data.Repo.FullName,
		"Today",
		healthClass, data.HealthScore,
		busClass, data.BusFactor,
		data.MaturityLevel,
		data.Repo.StargazersCount,
		data.Repo.ForksCount,
		len(data.Contributors),
		len(data.Commits),
		data.BusRisk,
		summary,
	)

	_, err = file.WriteString(finalHTML)
	return err
}
