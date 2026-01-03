# Repo-lyzer CLI Improvements - Developer Guide

## What's Changed

### New Files
1. **internal/ui/settings.go** - Settings management system
2. **internal/ui/history.go** - Analysis history tracking
3. **internal/ui/help.go** - Comprehensive help system

### Modified Files
1. **internal/ui/menu.go** - Complete redesign with hierarchical menus
2. **internal/ui/styles.go** - Added 10+ new style definitions
3. **internal/ui/app.go** - Integrated all new features, new state machine
4. **internal/ui/export.go** - Added CSV and HTML export formats
5. **internal/ui/dashboard.go** - Enhanced UI and export menu

## Architecture

### State Machine
The application now uses 8 states (previously 4):

```
stateMenu          → Main hierarchical menu
stateAnalysisMenu  → Analysis type selection
stateInput         → Repository input
stateLoading       → Analysis in progress
stateDashboard     → Results display
stateSettings      → Settings management
stateHelp          → Help system
stateHistory       → Recent analyses
stateCompareInput  → Compare repositories
```

### Menu Hierarchy
```
Main Menu (MenuLevel)
├── Analysis Options (AnalysisMenu)
│   ├── Quick Analysis
│   ├── Detailed Analysis
│   ├── Custom Report
│   └── Back
├── Compare Repositories
├── Settings (SettingsMenu)
│   ├── Theme
│   ├── GitHub Token
│   ├── Export Preferences
│   ├── Reset Defaults
│   └── Back
├── Help (HelpMenu)
│   ├── Shortcuts
│   ├── FAQ
│   ├── Getting Started
│   ├── Documentation
│   └── Back
├── Recent Analyses
└── Exit
```

## Key Components

### Settings System
- **Location**: `~/.repo-lyzer/config.json`
- **Structure**: `AppSettings` struct
- **Features**: 
  - Theme selection
  - GitHub token storage
  - Export preferences
  - Notification settings

### History Tracking
- **Location**: `~/.repo-lyzer/history.json`
- **Structure**: `AnalysisHistory` slice
- **Features**:
  - Auto-tracking of analyses
  - Last 20 entries retained
  - Quick re-run capability

### Help System
- **Topics**: 4 main help categories
- **Accessibility**: Press 'h' or '?' at any menu
- **Content**: Shortcuts, FAQ, guides, docs

### Export Options
- **Formats**: JSON, Markdown, CSV, HTML
- **Selection**: Interactive menu with arrow keys
- **Feedback**: Real-time status messages

## Integration Points

### For Adding New Features

#### 1. Add New Menu Item
In `menu.go`:
```go
// Add to getMainMenuChoices()
"🆕 New Feature",

// Handle in handleMainMenuSelection()
case 5: // New Feature
    m.SelectedVal = "new_feature"
    m.Done = true
```

In `app.go`:
```go
// Add new state
stateNewFeature sessionState = iota

// Handle in MainModel.Update()
case stateNewFeature:
    // Handle your feature
```

#### 2. Add New Export Format
In `export.go`:
```go
func ExportNewFormat(data AnalysisResult, filename string) error {
    // Your implementation
}
```

In `dashboard.go`:
```go
// Add to exportOptions in NewDashboardModel()
"🆕 Export as NewFormat",

// Handle in handleExportSelection()
case 4: // New Format
    return m, func() tea.Msg {
        err := ExportNewFormat(m.data, "analysis.fmt")
        return exportMsg{err, "Exported to analysis.fmt"}
    }
```

#### 3. Add New Setting
In `settings.go`:
```go
// Add field to AppSettings
type AppSettings struct {
    NewSetting bool `json:"new_setting"`
}

// Add to UI in NewSettingsModel()
"🆕 New Setting: " + boolToStr(settings.NewSetting),

// Handle in handleSettingSelection()
case 7: // New Setting
    m.settings.NewSetting = !m.settings.NewSetting
```

## Testing the Changes

### Manual Testing
```bash
go run main.go
```

Navigate through:
1. Main menu
2. Sub-menus
3. Settings
4. Help system
5. Recent analyses
6. Analysis and export

### Testing Checklist
- [ ] All menu items are accessible
- [ ] Navigation works (arrows, enter, escape)
- [ ] Settings save and load
- [ ] History tracking works
- [ ] All export formats work
- [ ] Help displays correctly
- [ ] Error handling shows helpful messages

## Code Style Guidelines

### Menu Functions
```go
func getXMenuChoices() []string {
    return []string{
        "Emoji Description",
        "Emoji Description",
        // ... more items
    }
}

func (m *EnhancedMenuModel) handleXMenuSelection() {
    switch m.cursor {
    case 0:
        // Handle first item
    case 1:
        // Handle second item
    // ...
    }
}
```

### Tea Model Pattern
```go
type MyModel struct {
    field1 Type
    field2 Type
}

func NewMyModel() MyModel {
    return MyModel{
        // Initialize fields
    }
}

func (m MyModel) Init() tea.Cmd { return nil }

func (m MyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle messages
    return m, nil
}

func (m MyModel) View() string {
    // Render UI
    return ""
}
```

## Performance Notes

1. **Settings**: Loaded once at startup, cached
2. **History**: Limited to 20 entries to prevent bloat
3. **Export**: Async operations with spinner
4. **UI Rendering**: Uses lipgloss for efficient terminal rendering

## Dependencies

No new external dependencies added. Uses existing:
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- Standard library JSON for config/history

## Debugging

### Enable Debug Output
Modify `app.go`:
```go
func (m MainModel) View() string {
    debugInfo := fmt.Sprintf("State: %d, Menu Level: %d", m.state, m.menu.level)
    // Add to output
}
```

### Check Config File
```bash
cat ~/.repo-lyzer/config.json
cat ~/.repo-lyzer/history.json
```

### Test Settings I/O
In `settings.go`:
```go
func TestSettings() {
    s := DefaultSettings
    SaveSettings(s)
    loaded, _ := LoadSettings()
    // Verify
}
```

## Common Issues & Solutions

### Settings Not Saving
- Check directory permissions: `~/.repo-lyzer/`
- Verify JSON is valid
- Check disk space

### History Not Appearing
- Ensure at least one analysis completed
- Check history file exists: `~/.repo-lyzer/history.json`
- Verify file is readable

### Menu Not Responding
- Check terminal width/height
- Verify tea.Msg handling in Update()
- Check for infinite loops in navigation

## Future Considerations

1. **Persistence**: Consider database instead of JSON files
2. **Performance**: Cache large history datasets
3. **Extensibility**: Plugin system for custom analyzers
4. **Internationalization**: Multi-language support
5. **Cloud Sync**: Sync across devices

## Resources

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI Framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Go Documentation](https://golang.org/doc/)

## Questions?

Refer to:
1. CLI_IMPROVEMENTS.md - User perspective
2. Code comments in .go files
3. Issue tracker for known problems
