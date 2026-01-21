package output

import (
	"fmt"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/charmbracelet/lipgloss"
)

var (
	certTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD700")).
		Align(lipgloss.Center).
		Margin(1, 0)

	certSectionStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF87")).
		MarginTop(1)

	certKeyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00E5FF")).
		Bold(true)

	certValueStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	scoreStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFB000")).
		Bold(true)

	gradeStyle = lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center)
)

// PrintCertificate displays a formatted repository certificate
func PrintCertificate(cert *analyzer.CertificateData) {
	// Title
	fmt.Println(certTitleStyle.Render("🏆 REPOSITORY CERTIFICATE 🏆"))
	fmt.Printf("Repository: %s/%s\n", cert.Owner, cert.RepoName)
	fmt.Println(strings.Repeat("=", 60))

	// Repository Information
	fmt.Println(certSectionStyle.Render("📋 Repository Information"))
	fmt.Printf("%s: %s\n", certKeyStyle.Render("Description"), certValueStyle.Render(cert.Description))
	fmt.Printf("%s: %d\n", certKeyStyle.Render("Stars"), certValueStyle.Render(fmt.Sprintf("%d", cert.Stars)))
	fmt.Printf("%s: %d\n", certKeyStyle.Render("Forks"), certValueStyle.Render(fmt.Sprintf("%d", cert.Forks)))
	fmt.Printf("%s: %d\n", certKeyStyle.Render("Open Issues"), certValueStyle.Render(fmt.Sprintf("%d", cert.OpenIssues)))
	fmt.Printf("%s: %s\n", certKeyStyle.Render("Created"), certValueStyle.Render(cert.CreatedAt))
	fmt.Printf("%s: %s\n", certKeyStyle.Render("Last Updated"), certValueStyle.Render(cert.UpdatedAt))
	fmt.Printf("%s: %s (%d languages)\n", certKeyStyle.Render("Primary Language"), certValueStyle.Render(cert.PrimaryLanguage), cert.LanguageCount)

	// Scores
	fmt.Println(certSectionStyle.Render("📊 Scores & Metrics"))
	fmt.Printf("%s: %s/100\n", certKeyStyle.Render("Health Score"), scoreStyle.Render(fmt.Sprintf("%d", cert.HealthScore)))
	fmt.Printf("%s: %s/100 (%s)\n", certKeyStyle.Render("Maturity Score"), scoreStyle.Render(fmt.Sprintf("%d", cert.MaturityScore)), certValueStyle.Render(cert.MaturityLevel))
	fmt.Printf("%s: %s (%s)\n", certKeyStyle.Render("Bus Factor"), scoreStyle.Render(fmt.Sprintf("%d", cert.BusFactor)), certValueStyle.Render(cert.BusRisk))
	fmt.Printf("%s: %d (%s)\n", certKeyStyle.Render("Commits (Last Year)"), certValueStyle.Render(fmt.Sprintf("%d", cert.CommitsLastYear)), certValueStyle.Render(cert.ActivityLevel))
	fmt.Printf("%s: %d\n", certKeyStyle.Render("Contributors"), certValueStyle.Render(fmt.Sprintf("%d", cert.Contributors)))

	// Overall Assessment
	fmt.Println(certSectionStyle.Render("🎯 Overall Assessment"))
	fmt.Printf("%s: %s\n", certKeyStyle.Render("Overall Score"), scoreStyle.Render(fmt.Sprintf("%d/100", cert.OverallScore)))
	fmt.Printf("%s: ", certKeyStyle.Render("Grade"))
	switch cert.Grade {
	case "A+", "A":
		fmt.Println(gradeStyle.Foreground(lipgloss.Color("#00FF00")).Render(cert.Grade))
	case "B+", "B":
		fmt.Println(gradeStyle.Foreground(lipgloss.Color("#FFFF00")).Render(cert.Grade))
	case "C+", "C":
		fmt.Println(gradeStyle.Foreground(lipgloss.Color("#FFA500")).Render(cert.Grade))
	default:
		fmt.Println(gradeStyle.Foreground(lipgloss.Color("#FF0000")).Render(cert.Grade))
	}

	// Potential Uses
	fmt.Println(certSectionStyle.Render("💡 Potential Uses"))
	for i, use := range cert.Uses {
		fmt.Printf("%d. %s\n", i+1, certValueStyle.Render(use))
	}

	fmt.Println(strings.Repeat("=", 60))
}
