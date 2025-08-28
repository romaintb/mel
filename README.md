# Next-Gen TUI Email Client

A modern terminal-based email client built in Go that combines the efficiency of CLI tools with the polish of contemporary TUI applications.

## Features

- **Resource Efficient**: Lightweight alternative to memory-heavy webmail clients
- **Modern TUI**: Inspired by lazygit and neovim aesthetics with mouse + keyboard support
- **Conversation Threading**: Gmail-style email threading as first-class experience
- **Powerful Search**: Telescope-style fuzzy search with live preview
- **Zero Config**: Auto-detects existing maildir setups from mbsync/offlineimap
- **Neovim-inspired Interface**: Modal operations with visual discoverability
- **Gmail-style Threading**: Conversation-first email paradigm
- **Resource Conscious**: Sub-100MB memory footprint, instant startup
- **External Tool Integration**: Works with mbsync, notmuch, and msmtp
- **Keyboard-First**: All operations accessible via shortcuts
- **Mouse Support**: Enhanced interaction for discoverability

## Architecture

Follows the aerc philosophy - delegates email sync/send to external tools (mbsync, msmtp, notmuch) and focuses on delivering an exceptional user experience.

## Documentation

- **[Product Requirements Document](docs/PRD.md)** - Complete product vision and technical requirements
- **Setup Guide** - Coming soon
- **Architecture Overview** - Coming soon

## Quick Start

### Prerequisites

- Go 1.24 or later
- External email tools (mbsync, notmuch, msmtp)

### Installation

```bash
# Clone the repository
git clone https://github.com/romaintb/mel.git
cd mel

# Build the application
make build

# Run the application
./bin/mel
```

### Development Setup

```bash
# Install development tools
make install-tools

# Run tests
make test

# Format code
make fmt

# Run linter
make lint

# Run all checks
make check
```

## Keybindings

### Normal Mode
- `h` - Focus sidebar
- `l` - Focus thread list
- `i` - Enter insert mode
- `v` - Enter visual mode
- `/` - Enter search mode
- `q` - Quit

### Thread Navigation
- `j/k` - Navigate threads
- `gg/G` - Go to first/last thread
- `enter` - Select thread

### Search
- `<leader>fg` - Content search
- `<leader>fs` - Sender search
- `<leader>fe` - Global email search

## Architecture

Mel follows a delegation model where external tools handle email operations:

- **mbsync/offlineimap**: IMAP synchronization
- **notmuch**: Email indexing and threading
- **msmtp**: SMTP sending

The application focuses on providing an excellent user experience while delegating the complex email handling to proven tools.

## Project Structure

```
mel/
├── cmd/mel/          # Main application entry point
├── internal/         # Private application code
│   ├── app/         # Main application logic
│   ├── config/      # Configuration management
│   └── ui/          # TUI components
├── pkg/             # Potentially reusable packages
├── docs/            # Documentation
└── scripts/         # Build and utility scripts
```

## Configuration

Configuration is automatically loaded from `~/.config/mel/config.yaml`. The application will create a default configuration if none exists.

### Example Configuration

```yaml
email:
  maildir: ~/Mail
  default_account: ""
  auto_sync_interval: 300

ui:
  theme:
    color_scheme: "auto"
    show_unread_indicators: true
    show_sync_status: true
  keybindings:
    leader: " "

external_tools:
  mbsync: "mbsync"
  notmuch: "notmuch"
  msmtp: "msmtp"
```

## Development

### Code Style

- Follow Go best practices and idioms
- Use 120 character line width
- Prefer clear code over comments
- Keep models and controllers thin
- Use services for business logic

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make bench
```

### Linting

```bash
# Run linter
make lint

# Run go vet
make vet
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by [lazygit](https://github.com/jesseduffield/lazygit)
- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
