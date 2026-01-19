# Fix Errors in Branch

## Tasks
- [x] Fix invalid Go version in go.mod (1.24.4 → 1.21)
- [x] Fix incorrect quick access key mappings in menu.go
  - [x] "c" (Compare): cursor 1 → 2
  - [x] "h" (History): cursor 2 → 3
  - [x] "s" (Settings): cursor 4 → 5
- [x] Fix TODO for export status messages in app.go
  - [x] Add success messages for JSON export
  - [x] Add success messages for Markdown export

## Status
- All identified errors have been fixed
- Branch should now build and run correctly
