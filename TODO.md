# Bug Fix: No output displayed in certain situations

## Description
Repo-lyzer does not display any output when analyzing repositories with no issues or when comparing repositories with no differences.

## Tasks
- [ ] Modify output functions in `internal/output/` to display messages when no data is available
- [ ] Update `PrintCommitActivity` to show "No commit activity found" when data is empty
- [ ] Update `PrintLanguages` to show "No language data available" when map is empty
- [ ] Update compare command to detect identical repositories and show "No differences found" message
- [ ] Test the changes with repositories having no data

## Files to Edit
- internal/output/charts.go
- internal/output/json.go
- cmd/compare.go
- internal/ui/app.go (for compare result view)
