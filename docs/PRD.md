# Next-Gen TUI Email Client - PRD

## Overview

### Product Vision
A modern terminal-based email client that combines the efficiency of CLI tools with the polish of contemporary TUI applications. Delivers the resource efficiency and keyboard-driven speed that power users crave, while maintaining the familiar interaction patterns of modern development tools.

### Problem Statement
- **Resource Inefficiency**: Webmail clients consume hundreds of megabytes of RAM for what is fundamentally text processing
- **Outdated UX**: Existing TUI email clients feel antiquated with clunky interfaces and arcane keyboard shortcuts
- **Threading Gap**: Most terminal email clients handle email-by-email workflows poorly, lacking modern conversation threading

### Target User
Primary: Developers and power users who live in terminal environments (neovim users, CLI-first workflows) but currently use resource-heavy webmail clients for convenience.

## Core Principles

### Design Philosophy
- **Familiar Efficiency**: Keyboard shortcuts and interaction patterns that feel natural to neovim/lazygit users
- **Hybrid Interaction**: Mouse support for discoverability, keyboard shortcuts for speed
- **Conversation-First**: Gmail-style threading as the default email paradigm
- **Resource Conscious**: Sub-100MB memory footprint, instant startup, snappy interactions

### Technical Architecture
- **Delegation Model**: Follow aerc's approach - external tools handle SMTP/IMAP/sync, client focuses on UX
- **Go Implementation**: Leverage Go's excellent concurrency model and TUI ecosystem
- **Modal Interface**: Neovim-inspired modal operations with visual discoverability

## Core Features

### 1. Navigation & Layout

#### Sidebar System
- **Left Sidebar**: Collapsible account/folder tree (nvim-tree style)
  - Visual indicators for unread counts, sync status
  - Keyboard navigation (`hjkl`) and mouse interaction
  - Toggle with `<leader>e`

#### Main Panel Layout
- **Thread List**: Gmail-style conversation view
  - Subject, participants summary, timestamp, unread indicators
  - Thread-level actions (archive, delete, star) by default
- **Thread View**: Continuous conversation flow
  - Latest messages expanded, older messages collapsed
  - Inline expansion without losing context

#### Modal Operations
- **Normal Mode**: Thread/email navigation, bulk operations
- **Insert Mode**: Composition and reply
- **Visual Mode**: Multi-selection for bulk actions
- Context-sensitive shortcuts with visual hints

### 2. Search & Discovery

#### Telescope Integration
- **Two-pane search**: Results list + live preview
- **Content Search** (`<leader>fg`): Full-text search across all emails with highlighted matches
- **Sender Search** (`<leader>fs`): Find all emails from specific people/addresses
- **Global Email Search** (`<leader>fe`): Fuzzy search across subjects, senders, dates

#### Search Capabilities
- Live preview of email content while browsing results
- Incremental search with instant feedback
- Context snippets showing match locations
- Quick actions from search results (delete, archive, etc.)

### 3. Threading & Conversations

#### Gmail-Style Threading
- **Single Conversation Flow**: All emails displayed as continuous discussion
- **Smart Collapsing**: Show only latest 2-3 messages by default
- **Expand on Demand**: "Show X earlier messages" functionality
- **Thread Actions**: Operations apply to entire conversation by default

#### Thread Navigation
- `j/k` between individual emails within thread
- `n/p` for next/previous thread
- `o` to expand/collapse individual messages
- Smart auto-scroll to first unread message

### 4. Composition

#### Compose Interface
- **Context-Aware**: Show thread history above compose area for replies
- **Modal Editing**: Enter insert mode for writing, normal mode for operations
- **Quick Actions**: Fast reply, forward, compose new

## Technical Requirements

### External Tool Integration
- **Email Sync**: `mbsync`/`offlineimap` for IMAP synchronization
- **Sending**: `msmtp`/`sendmail` for SMTP
- **Indexing**: `notmuch` for search and threading support
- **Process Management**: Go's goroutines for concurrent external tool coordination

### Performance Targets
- **Memory Usage**: < 100MB runtime footprint
- **Startup Time**: < 500ms cold start
- **Response Time**: < 100ms for all UI interactions
- **Search Speed**: Instant results for indexed content

### Platform Support
- **Primary**: Linux, macOS (Unix-like systems)
- **Dependencies**: Standard Unix tools, configurable external programs

## User Experience Goals

### Interaction Model
- **Keyboard-First**: All operations accessible via shortcuts
- **Mouse-Enhanced**: Click for selection, scroll for navigation
- **Progressive Disclosure**: Visual hints for available actions
- **Contextual Help**: Bottom status bar showing relevant shortcuts

### Familiarity Targets
- **Neovim Users**: Modal editing, familiar movement patterns (`hjkl`, `gg/G`, `/` search)
- **Lazygit Aesthetic**: Clean panels, contextual shortcuts, snappy interactions
- **Modern CLI Tools**: Polish level of `btop`, `bat`, `ripgrep`

## Success Metrics

### User Adoption (Personal Use)
- **Zero-Config Goal**: Launch and immediately see emails from existing sync setups
- **Primary Success**: Daily driver replacement for webmail
- **Workflow Speed**: Measurably faster common operations (archive, search, compose)
- **Resource Efficiency**: Demonstrable memory/CPU savings vs. webmail

### Technical Milestones
- **MVP**: Core threading, search, compose functionality
- **Daily Use Ready**: Stable performance with real email volumes
- **Feature Complete**: Full feature parity with primary webmail workflows

## Implementation Phases

### Phase 1: Core Foundation (2-3 months)
- Basic TUI framework with bubbletea/tview
- External tool integration (mbsync, notmuch)
- Simple email list and reading interface
- Basic threading support

### Phase 2: Search & Navigation (1-2 months)
- Telescope-style search implementation
- Advanced keyboard navigation
- Thread conversation view
- Sidebar and panel management

### Phase 3: Polish & Composition (1-2 months)
- Compose interface with context
- Mouse support integration
- Visual polish and status indicators
- Configuration system

### Phase 4: Advanced Features (Ongoing)
- Smart threading improvements
- Performance optimizations
- Extended external tool support
- Plugin/scripting capabilities

## Configuration Philosophy

### Approach
- **Auto-Detection**: Automatically discover emails in standard sync tool directories (`~/Mail`, `~/.mail`, etc.)
- **Sensible Defaults**: Works well out of the box with common mbsync/offlineimap setups
- **Progressive Configuration**: Simple config for basic use, powerful options for customization
- **External Tool Agnostic**: Support multiple backends for sync/send/index operations

### Configuration Areas
- **Account Discovery**: Auto-detect maildir structures from sync tools
- **External Tool Integration**: Configurable paths for mbsync, notmuch, msmtp
- **Hardcoded Keybindings**: Neovim-style shortcuts built-in, non-remappable for consistency
- **Theme System**: Color schemes and visual customization (planned for later phases)
- **Fixed Layout**: No layout customization initially - focus on perfecting default experience
- Search and threading behavior

## Risks & Considerations

### Technical Challenges
- **Email Complexity**: Handling various encodings, formats, edge cases
- **Threading Logic**: Robust conversation detection across different email clients
- **Performance**: Maintaining responsiveness with large mailboxes

### User Adoption
- **Learning Curve**: Modal interface may require adjustment period
- **External Dependencies**: Users need to configure and maintain external tools
- **Feature Gaps**: May lack some webmail conveniences initially

### Mitigation Strategies
- Comprehensive documentation and setup guides
- Gradual feature rollout focusing on core workflows first
- Community feedback integration throughout development
- Fallback mechanisms for edge cases