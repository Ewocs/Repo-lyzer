package ui

import "github.com/agnivo988/Repo-lyzer/internal/github"

type AnalysisResult struct {
	Repo          *github.Repo
	Commits       []github.Commit
	Contributors  []github.Contributor
	Languages     map[string]int
	HealthScore   int
	BusFactor     int
	BusRisk       string
	MaturityScore int
	MaturityLevel string
}
