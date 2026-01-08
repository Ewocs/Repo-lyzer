// Package analyzer provides security vulnerability scanning for dependencies.
// This file implements vulnerability detection using the OSV (Open Source Vulnerabilities) database.
package analyzer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Vulnerability represents a security vulnerability found in a dependency
type Vulnerability struct {
	ID          string   `json:"id"`           // CVE or GHSA ID
	Summary     string   `json:"summary"`      // Brief description
	Severity    string   `json:"severity"`     // CRITICAL, HIGH, MEDIUM, LOW
	Package     string   `json:"package"`      // Affected package name
	Version     string   `json:"version"`      // Affected version
	FixedIn     string   `json:"fixed_in"`     // Version that fixes the issue
	References  []string `json:"references"`   // Links to advisories
	PublishedAt string   `json:"published_at"` // When vulnerability was published
}

// SecurityScanResult holds the complete security scan results
type SecurityScanResult struct {
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	TotalCount      int             `json:"total_count"`
	CriticalCount   int             `json:"critical_count"`
	HighCount       int             `json:"high_count"`
	MediumCount     int             `json:"medium_count"`
	LowCount        int             `json:"low_count"`
	ScannedPackages int             `json:"scanned_packages"`
	ScanTime        time.Time       `json:"scan_time"`
	SecurityScore   int             `json:"security_score"` // 0-100, higher is better
}

// OSV API structures
type osvQuery struct {
	Package struct {
		Name      string `json:"name"`
		Ecosystem string `json:"ecosystem"`
	} `json:"package"`
	Version string `json:"version,omitempty"`
}

type osvResponse struct {
	Vulns []osvVuln `json:"vulns"`
}

type osvVuln struct {
	ID       string `json:"id"`
	Summary  string `json:"summary"`
	Details  string `json:"details"`
	Severity []struct {
		Type  string `json:"type"`
		Score string `json:"score"`
	} `json:"severity"`
	Affected []struct {
		Package struct {
			Name      string `json:"name"`
			Ecosystem string `json:"ecosystem"`
		} `json:"package"`
		Ranges []struct {
			Type   string `json:"type"`
			Events []struct {
				Introduced string `json:"introduced,omitempty"`
				Fixed      string `json:"fixed,omitempty"`
			} `json:"events"`
		} `json:"ranges"`
	} `json:"affected"`
	References []struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"references"`
	Published string `json:"published"`
}

// ScanDependencies scans dependencies for known vulnerabilities using OSV database
func ScanDependencies(deps *DependencyAnalysis) (*SecurityScanResult, error) {
	if deps == nil || len(deps.Files) == 0 {
		return &SecurityScanResult{
			Vulnerabilities: []Vulnerability{},
			ScanTime:        time.Now(),
			SecurityScore:   100,
		}, nil
	}

	result := &SecurityScanResult{
		Vulnerabilities: []Vulnerability{},
		ScanTime:        time.Now(),
	}

	client := &http.Client{Timeout: 10 * time.Second}

	for _, file := range deps.Files {
		ecosystem := mapFileTypeToEcosystem(file.FileType)
		if ecosystem == "" {
			continue
		}

		for _, dep := range file.Dependencies {
			result.ScannedPackages++

			// Query OSV API
			vulns, err := queryOSV(client, dep.Name, dep.Version, ecosystem)
			if err != nil {
				continue // Skip on error, don't fail entire scan
			}

			for _, v := range vulns {
				vuln := convertOSVVuln(v, dep.Name, dep.Version)
				result.Vulnerabilities = append(result.Vulnerabilities, vuln)

				// Count by severity
				switch vuln.Severity {
				case "CRITICAL":
					result.CriticalCount++
				case "HIGH":
					result.HighCount++
				case "MEDIUM":
					result.MediumCount++
				case "LOW":
					result.LowCount++
				}
			}
		}
	}

	result.TotalCount = len(result.Vulnerabilities)
	result.SecurityScore = calculateSecurityScore(result)

	return result, nil
}

// mapFileTypeToEcosystem maps our file types to OSV ecosystem names
func mapFileTypeToEcosystem(fileType string) string {
	ecosystems := map[string]string{
		"npm":    "npm",
		"go":     "Go",
		"python": "PyPI",
		"rust":   "crates.io",
		"ruby":   "RubyGems",
	}
	return ecosystems[fileType]
}

// queryOSV queries the OSV API for vulnerabilities
func queryOSV(client *http.Client, packageName, version, ecosystem string) ([]osvVuln, error) {
	query := osvQuery{}
	query.Package.Name = packageName
	query.Package.Ecosystem = ecosystem

	// Clean version string
	cleanVer := cleanVersionForOSV(version)
	if cleanVer != "" && cleanVer != "*" {
		query.Version = cleanVer
	}

	jsonData, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(
		"https://api.osv.dev/v1/query",
		"application/json",
		strings.NewReader(string(jsonData)),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OSV API error: %d", resp.StatusCode)
	}

	var osvResp osvResponse
	if err := json.NewDecoder(resp.Body).Decode(&osvResp); err != nil {
		return nil, err
	}

	return osvResp.Vulns, nil
}

// cleanVersionForOSV removes version prefixes for OSV API
func cleanVersionForOSV(version string) string {
	version = strings.TrimSpace(version)
	// Remove common prefixes
	version = strings.TrimPrefix(version, "^")
	version = strings.TrimPrefix(version, "~")
	version = strings.TrimPrefix(version, ">=")
	version = strings.TrimPrefix(version, "<=")
	version = strings.TrimPrefix(version, ">")
	version = strings.TrimPrefix(version, "<")
	version = strings.TrimPrefix(version, "=")
	version = strings.TrimPrefix(version, "v")
	return version
}

// convertOSVVuln converts OSV vulnerability to our format
func convertOSVVuln(osv osvVuln, pkgName, pkgVersion string) Vulnerability {
	vuln := Vulnerability{
		ID:          osv.ID,
		Summary:     osv.Summary,
		Package:     pkgName,
		Version:     pkgVersion,
		PublishedAt: osv.Published,
	}

	// Extract severity
	vuln.Severity = extractSeverity(osv)

	// Extract fixed version
	for _, affected := range osv.Affected {
		for _, r := range affected.Ranges {
			for _, event := range r.Events {
				if event.Fixed != "" {
					vuln.FixedIn = event.Fixed
					break
				}
			}
		}
	}

	// Extract references
	for _, ref := range osv.References {
		if ref.URL != "" {
			vuln.References = append(vuln.References, ref.URL)
		}
	}

	// Limit references to 3
	if len(vuln.References) > 3 {
		vuln.References = vuln.References[:3]
	}

	return vuln
}

// extractSeverity extracts severity from OSV vulnerability
func extractSeverity(osv osvVuln) string {
	for _, sev := range osv.Severity {
		if sev.Type == "CVSS_V3" {
			score := parseCVSSScore(sev.Score)
			if score >= 9.0 {
				return "CRITICAL"
			} else if score >= 7.0 {
				return "HIGH"
			} else if score >= 4.0 {
				return "MEDIUM"
			}
			return "LOW"
		}
	}

	// Default based on ID prefix
	if strings.HasPrefix(osv.ID, "GHSA") {
		return "MEDIUM" // Default for GitHub advisories
	}
	return "UNKNOWN"
}

// parseCVSSScore parses CVSS score from string
func parseCVSSScore(score string) float64 {
	var s float64
	fmt.Sscanf(score, "%f", &s)
	return s
}

// calculateSecurityScore calculates overall security score (0-100)
func calculateSecurityScore(result *SecurityScanResult) int {
	if result.ScannedPackages == 0 {
		return 100
	}

	// Deduct points based on vulnerabilities
	score := 100
	score -= result.CriticalCount * 25
	score -= result.HighCount * 15
	score -= result.MediumCount * 5
	score -= result.LowCount * 2

	if score < 0 {
		score = 0
	}
	return score
}

// GetSeverityEmoji returns emoji for severity level
func GetSeverityEmoji(severity string) string {
	switch severity {
	case "CRITICAL":
		return "🔴"
	case "HIGH":
		return "🟠"
	case "MEDIUM":
		return "🟡"
	case "LOW":
		return "🟢"
	default:
		return "⚪"
	}
}

// GetSecurityGrade returns letter grade based on score
func GetSecurityGrade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}
