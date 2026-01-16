# Fix Error Message for Invalid Repository URLs

## Tasks
- [x] Update error messages in cmd/analyze.go, internal/ui/app.go for invalid repository URL format
- [x] Ensure consistent error messaging across CLI and UI interfaces

## Status
- Analysis complete: Error messages were too generic and didn't clearly indicate invalid URL format
- Plan approved: Update error messages to be more descriptive about valid formats
- Implementation complete: Error messages now clearly state "invalid repository URL: must be in owner/repo format or a valid GitHub URL"
