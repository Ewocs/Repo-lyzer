# Repo-lyzer CLI - Quick Reference Card

## 🚀 Quick Start

### Launch
```bash
go run main.go
./repo-lyzer
```

### First Time
1. Press 'h' for help
2. Navigate with ↑/↓
3. Press Enter to select

---

## 📋 Main Menu Structure

```
┌─────────────────────────────────────┐
│        REPO-LYZER MAIN MENU         │
├─────────────────────────────────────┤
│                                     │
│  📊 Analyze Repository              │
│  ├─ 🔍 Quick Analysis               │
│  ├─ 📈 Detailed Analysis            │
│  └─ 📊 Custom Report                │
│                                     │
│  🔄 Compare Repositories            │
│                                     │
│  ⚙️  Settings                       │
│  ├─ 🎨 Theme Selection              │
│  ├─ 🔐 GitHub Token                 │
│  ├─ 📁 Export Preferences           │
│  └─ 🔄 Reset Defaults               │
│                                     │
│  ℹ️  Help                           │
│  ├─ ⌨️  Shortcuts                   │
│  ├─ ❓ FAQ                          │
│  ├─ 📖 Getting Started              │
│  └─ 🔗 Documentation                │
│                                     │
│  📜 Recent Analyses                 │
│  ❌ Exit                            │
│                                     │
└─────────────────────────────────────┘
```

---

## ⌨️ Keyboard Shortcuts

| Key | Action | Where |
|-----|--------|-------|
| ↑ | Move up | All menus |
| ↓ | Move down | All menus |
| k | Move up (vim) | All menus |
| j | Move down (vim) | All menus |
| Enter | Select | All menus |
| ESC | Back | Submenus |
| b | Back (alt) | Submenus |
| q | Quit | Main menu |
| h | Help toggle | All menus |
| ? | Help toggle | All menus |
| e | Export | Dashboard |
| Ctrl+C | Force quit | Anywhere |

---

## 📊 Workflow Examples

### Example 1: Analyze a Repository

```
1. Launch app
2. ↓ Select "Analyze Repository"
3. Press Enter
4. ↓↓ Select "Quick Analysis" (or Detailed/Custom)
5. Press Enter
6. Type: owner/repo (e.g., golang/go)
7. Press Enter
8. Wait for analysis
9. View results
10. Press 'e' to export
11. Select format (JSON/Markdown/CSV/HTML)
12. Press 'q' to return to main menu
```

### Example 2: Configure Settings

```
1. Launch app
2. ↓↓ Select "Settings"
3. Press Enter
4. ↓ Select setting to change
5. Press Enter to toggle/cycle
6. View changes immediately
7. Press ESC to return
8. Settings auto-save
```

### Example 3: Get Help

```
1. At any menu, press 'h'
2. Help text appears below menu
3. Read the help content
4. Press 'h' again to hide
5. Or go to Help menu for full topics
```

### Example 4: Use Recent History

```
1. Select "Recent Analyses"
2. Select previous repo
3. Analysis re-runs automatically
4. View results faster
```

---

## 🎨 Color Guide

| Color | Meaning | Example |
|-------|---------|---------|
| 🔵 Cyan | Important | Menu title, metrics |
| 🟢 Green | Success | Selected item, ✅ |
| ⚪ White | Normal | Menu options |
| 🟡 Gold | Input | Repository name |
| ⚫ Gray | Subtle | Help text, shortcuts |
| 🔴 Red | Error | ❌ Error messages |
| 🟠 Orange | Warning | ⚠️ Important notes |

---

## 📁 Files & Directories

### Configuration
```
~/.repo-lyzer/
├── config.json      (Settings)
└── history.json     (Analysis history)
```

### Settings Example
```json
{
  "theme": "dark",
  "github_token": "ghp_...",
  "default_export_type": "json",
  "default_export_path": "./exports",
  "show_progress_bars": true,
  "enable_notifications": true
}
```

### Export Output
```
./analysis.json       (Data export)
./analysis.md         (Report)
./analysis.csv        (Spreadsheet)
./analysis.html       (Web report)
```

---

## 💡 Tips & Tricks

### Power User Tips
1. Use vim keybindings (j/k) - faster navigation
2. Check recent analyses - skip re-entering names
3. Configure token once - runs faster
4. Bookmark help topics - quick reference
5. Export to HTML - share reports easily

### Keyboard Efficiency
1. Use ESC instead of navigating back
2. Use 'q' to quit from main menu
3. Use 'h' for contextual help
4. Combine keys: Select → Enter → Done

### Best Practices
1. Set GitHub token in Settings first
2. Start with Quick Analysis
3. Use Detailed for full insights
4. Export HTML for sharing
5. Check Recent Analyses for patterns

### Troubleshooting
1. **Settings not saving**: Check ~/.repo-lyzer/ permissions
2. **Analysis fails**: Verify GitHub token in Settings
3. **No history**: Run analysis to populate
4. **Export fails**: Check export path permissions
5. **Menu stuck**: Press Ctrl+C to restart

---

## 📊 Analysis Metrics Explained

### Health Score (0-100)
- **70+**: Good health ✅
- **50-70**: Moderate health ⚠️
- **<50**: Low health ❌

### Bus Factor (1-10+)
- **7+**: Low risk ✅
- **4-6**: Medium risk ⚠️
- **1-3**: High risk ❌

### Maturity Levels
- **Experimental**: New project
- **Development**: Active development
- **Stable**: Established, regular maintenance
- **Mature**: Long-running, well-maintained

---

## 🚀 Getting Help

### Built-in Help
- Press 'h' at any menu
- Select "Help" from main menu
- 4 help topics available

### Help Topics
1. **Shortcuts** - All keyboard shortcuts
2. **FAQ** - Frequently asked questions
3. **Getting Started** - Quick tutorial
4. **Documentation** - Full docs

### GitHub
- Issues: Report problems
- Discussions: Ask questions
- Pull Requests: Contribute

---

## 📈 Performance Tips

### Fast Analysis
1. Use Quick Analysis mode
2. Use Recent Analyses
3. Save GitHub token
4. Check README_changes for recent improvements

### Efficient Workflow
1. Batch analyze similar repos
2. Export once, review multiple
3. Use history for comparison
4. Configure once, reuse settings

---

## 🔧 Configuration Tips

### GitHub Token
1. Go to Settings → GitHub Token
2. Generate token on GitHub (Settings → Developer)
3. Paste token
4. Increases API rate limits
5. Auto-saved

### Theme Selection
1. Settings → Theme
2. Cycle through options
3. Takes effect immediately
4. Saved for next session

### Export Preferences
1. Settings → Export Preferences
2. Choose default format
3. Set export path
4. Used for quick exports

---

## ⚡ Quick Reference

### Menu Navigation
```
Main Menu → Submenu → Option → Result → Return
    ↓        ↓        ↓        ↓      ↑
  Enter    Enter    Enter   Display  ESC
```

### Analysis Flow
```
Analyze → Select Type → Enter Repo → Loading → Dashboard → Export
  ↓          ↓            ↓          ↓          ↓         ↓
Select     Choice       Input      Spinner    View      Format
```

### Settings Flow
```
Settings → Select Setting → Change Value → Auto-Save → Return
   ↓           ↓               ↓             ↓         ↓
Select      Navigate        Enter/Cycle   Feedback   ESC
```

---

## 📱 Terminal Requirements

- **Width**: 80+ characters recommended
- **Height**: 24+ lines recommended
- **Colors**: 256 color terminal
- **Encoding**: UTF-8 support

---

## 🎯 Common Tasks

### Task: Analyze golang/go
1. Select "Analyze Repository" → Enter
2. Select "Quick Analysis" → Enter
3. Type: golang/go
4. Press Enter
5. View results

### Task: Change Theme
1. Select "Settings" → Enter
2. Select "Theme" → Enter
3. Cycles through themes
4. Return with ESC

### Task: Export as HTML
1. After analysis, press 'e'
2. ↓↓↓ Select "Export as HTML"
3. Press Enter
4. File saved to ./analysis.html

### Task: View Recent Repos
1. Select "Recent Analyses" → Enter
2. View last 20 analyses
3. Select one → Press Enter
4. Auto re-runs analysis

---

## 🔐 Security Notes

- **Token Storage**: Local file (~/.repo-lyzer/)
- **History Storage**: Local file (~/.repo-lyzer/)
- **Network**: Only GitHub API calls
- **Privacy**: No data sent externally
- **Permissions**: Set restrictively (600)

---

## 📚 More Information

- **CLI_IMPROVEMENTS.md** - Full feature guide
- **DEVELOPER_GUIDE.md** - For developers
- **IMPLEMENTATION_SUMMARY.md** - Overview
- **VISUAL_CHANGES_GUIDE.md** - Visual reference

---

**Version**: 2.0 (with CLI improvements)
**Last Updated**: January 3, 2025
**Status**: Ready to Use ✅
