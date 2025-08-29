# Mel - Next-Gen TUI Email Client

A modern terminal-based email client that combines the efficiency of CLI tools with the polish of contemporary TUI applications. Built with Go and Bubble Tea, designed for developers and power users who live in terminal environments.

## 🎯 Product Vision

Mel delivers the resource efficiency and keyboard-driven speed that power users crave, while maintaining the familiar interaction patterns of modern development tools like neovim and lazygit.

## ✨ Key Features

### **Navigation & Layout**
- **Left Sidebar**: Collapsible account/folder tree with unread counts and sync status
- **Thread List**: Gmail-style conversation view with subject, participants, and timestamps
- **Thread View**: Continuous conversation flow with smart collapsing
- **Modal Interface**: Neovim-inspired modal operations (Normal/Insert/Visual/Search)

### **Search & Discovery**
- **Telescope Integration**: Two-pane search with results list and live preview
- **Content Search** (`<leader>fg`): Full-text search across all emails
- **Sender Search** (`<leader>fs`): Find emails from specific people/addresses
- **Global Email Search** (`<leader>fe`): Fuzzy search across subjects, senders, dates

### **Threading & Conversations**
- **Gmail-Style Threading**: Single conversation flow with smart collapsing
- **Thread Actions**: Archive, delete, star, mark read/unread
- **Smart Navigation**: `j/k` between emails, `n/p` between threads

### **External Tool Integration**
- **Email Sync**: `mbsync`/`offlineimap` for IMAP synchronization
- **Search & Indexing**: `notmuch` for search and threading support
- **Sending**: `msmtp`/`sendmail` for SMTP operations

## 🚀 Getting Started

### Prerequisites

Mel requires Go and these external tools to be installed and configured:

- **Go** ≥ 1.22 (to build; verify with `go version`)
- **mbsync** or **offlineimap** for email synchronization
- **notmuch** for search and threading
- **msmtp** or **sendmail** for sending emails
### Installation

```bash
# Clone the repository
git clone https://github.com/romaintb/mel.git
cd mel

# Build the application
make build

# Run mel
./cmd/mel/mel
```

### Configuration

Mel automatically detects your email setup and works with standard configurations:

- **Maildir**: `~/Mail` (configurable)
- **Config**: `~/.config/mel/config.yaml` (auto-generated with defaults)

#### **Icon Modes**

Mel supports two icon display modes:

- **ASCII Mode** (default): Uses single-character ASCII art icons for compatibility with all terminals
- **Emoji Mode**: Uses colorful emoji icons for better visual appeal

Toggle between modes with `<leader>i` or configure the default in your config file:

```yaml
ui:
  icon_mode: "ascii"  # or "emoji"
```

## ⌨️ Keybindings

### **Normal Mode**
- `h/l` - Navigate between sidebar and content
- `j/k` - Navigate threads
- `n/p` - Next/previous unread thread
- `o` - Expand/collapse thread
- `a` - Archive thread
- `d` - Delete thread
- `s` - Star/unstar thread
- `r` - Mark as read
- `u` - Mark as unread
- `e` - Toggle sidebar
- `i` - Enter insert mode
- `v` - Enter visual mode
- `/` - Enter search mode

### **Search Mode**
- `<leader>fg` - Content search
- `<leader>fs` - Sender search  
- `<leader>fe` - Global search
- `esc` - Exit search mode

### **Leader Key**
- `<space>` - Show available commands
- `<space>e` - Toggle sidebar
- `<space>i` - Toggle icon mode (emoji/ascii)

## 🏗️ Architecture

Mel follows a **delegation model** where external tools handle the heavy lifting:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│      Mel        │    │   External      │    │   Email        │
│   (TUI/UX)     │◄──►│     Tools       │◄──►│   Servers      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                       │
        │              ┌────────┴────────┐
        │              │                 │
        │           mbsync          notmuch
        │        (IMAP sync)      (search/index)
        │
        └──────────────┐
                       │
                   msmtp
                (SMTP send)
```

### **Core Components**
- **Email Manager**: Coordinates external tool operations
- **Search Service**: Implements Telescope-style search with relevance scoring
- **UI Components**: Modal interface with neovim-inspired interactions
- **Configuration**: Auto-detection with sensible defaults

## 📊 Performance Targets

- **Memory Usage**: < 100MB runtime footprint
- **Startup Time**: < 500ms cold start
- **Response Time**: < 100ms for all UI interactions
- **Search Speed**: Instant results for indexed content

## 🔧 Development

### **Building**
```bash
make build          # Build the application
make test           # Run tests
make clean          # Clean build artifacts
```

### **Project Structure**
```
internal/
├── app/           # Main application logic
├── config/        # Configuration management
├── email/         # Email data models and external tool integration
├── icons/         # Icon service with emoji/ASCII mode support
├── search/        # Search service with relevance scoring
└── ui/            # TUI components and modal interface
```

### **External Tool Integration**
- **mbsync**: Email synchronization via IMAP
- **notmuch**: Search, indexing, and threading
- **msmtp**: SMTP sending operations

Security note: Avoid storing SMTP/IMAP passwords in plain text. Prefer OAuth2 or OS keychain helpers (e.g., `pass`, GNOME Keyring, macOS Keychain) and restrict file permissions on config files.
## 🎨 Design Philosophy

- **Familiar Efficiency**: Keyboard shortcuts that feel natural to neovim/lazygit users
- **Hybrid Interaction**: Mouse support for discoverability, keyboard shortcuts for speed
- **Conversation-First**: Gmail-style threading as the default email paradigm
- **Resource Conscious**: Sub-100MB memory footprint, instant startup, snappy interactions

## 📈 Roadmap

### **Phase 1: Core Foundation** ✅
- [x] Basic TUI framework with Bubble Tea
- [x] External tool integration structure
- [x] Modal interface implementation
- [x] Email data models

### **Phase 2: Search & Navigation** 🚧
- [x] Search service architecture
- [x] Basic keybindings
- [ ] Complete search implementation
- [ ] Thread conversation view

### **Phase 3: Polish & Composition** 📋
- [ ] Compose interface with context
- [ ] Mouse support integration
- [ ] Visual polish and status indicators

### **Phase 4: Advanced Features** 📋
- [ ] Smart threading improvements
- [ ] Performance optimizations
- [ ] Extended external tool support

## 🤝 Contributing

Contributions are welcome! Please read our contributing guidelines and ensure your code follows our standards:

- All Go code must pass `gofmt`, `go vet`, and `staticcheck`
- Follow Go best practices and idioms
- Maintain the PRD-driven architecture
- Add tests for new functionality

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by the efficiency of `aerc` and the polish of modern TUI tools
- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for excellent TUI support
- Designed for the neovim/lazygit user experience
