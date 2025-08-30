package email

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// MailFolder represents a mail folder
type MailFolder struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	UnreadCount  int    `json:"unread_count"`
	MessageCount int    `json:"message_count"`
	IsSpecial    bool   `json:"is_special"` // Special folders like INBOX, Sent, etc.
}

// Thread represents a conversation thread
type Thread struct {
	ID            string     `json:"id"`
	Subject       string     `json:"subject"`
	Participants  []string   `json:"participants"`
	Timestamp     time.Time  `json:"timestamp"`
	UnreadCount   int        `json:"unread_count"`
	MessageCount  int        `json:"message_count"`
	LatestMessage *Message   `json:"latest_message"`
	Messages      []*Message `json:"messages"`
}

// Message represents an individual email message
type Message struct {
	ID        string    `json:"id"`
	ThreadID  string    `json:"thread_id"`
	From      string    `json:"from"`
	To        []string  `json:"to"`
	Cc        []string  `json:"cc"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	Timestamp time.Time `json:"timestamp"`
	Unread    bool      `json:"unread"`
	Starred   bool      `json:"starred"`
	Labels    []string  `json:"labels"`
}

// SearchResult represents a search result
type SearchResult struct {
	Threads []*Thread `json:"threads"`
	Query   string    `json:"query"`
	Total   int       `json:"total"`
}

// Manager handles email operations and external tool integration
type Manager struct {
	maildirPath string
	notmuchPath string
	mbsyncPath  string
	msmtpPath   string
}

// NewManager creates a new email manager
func NewManager(maildirPath, notmuchPath, mbsyncPath, msmtpPath string) *Manager {
	return &Manager{
		maildirPath: maildirPath,
		notmuchPath: notmuchPath,
		mbsyncPath:  mbsyncPath,
		msmtpPath:   msmtpPath,
	}
}

// SyncEmails synchronizes emails using mbsync
func (m *Manager) SyncEmails() error {
	cmd := exec.Command(m.mbsyncPath, "-a")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to sync emails: %w", err)
	}
	return nil
}

// SearchEmails searches emails using notmuch
func (m *Manager) SearchEmails(query string) (*SearchResult, error) {
	// Use notmuch search with JSON output
	cmd := exec.Command(m.notmuchPath, "search", "--format=json", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to search emails: %w", err)
	}

	// Parse notmuch JSON output and convert to our models
	// This is a simplified implementation - would need proper JSON parsing
	threads, err := m.parseNotmuchResults(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	return &SearchResult{
		Threads: threads,
		Query:   query,
		Total:   len(threads),
	}, nil
}

// GetThread retrieves a specific thread with all messages
func (m *Manager) GetThread(threadID string) (*Thread, error) {
	// Use notmuch show to get thread details
	cmd := exec.Command(m.notmuchPath, "show", "--format=json", threadID)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get thread: %w", err)
	}

	// Parse notmuch output and convert to Thread model
	thread, err := m.parseNotmuchThread(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse thread: %w", err)
	}

	return thread, nil
}

// MarkThreadRead marks all messages in a thread as read
func (m *Manager) MarkThreadRead(threadID string) error {
	cmd := exec.Command(m.notmuchPath, "tag", "-unread", fmt.Sprintf("thread:%s", threadID))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to mark thread as read: %w", err)
	}
	return nil
}

// ArchiveThread archives a thread (moves to archive folder)
func (m *Manager) ArchiveThread(threadID string) error {
	cmd := exec.Command(m.notmuchPath, "tag", "+archive", fmt.Sprintf("thread:%s", threadID))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to archive thread: %w", err)
	}
	return nil
}

// DeleteThread deletes a thread
func (m *Manager) DeleteThread(threadID string) error {
	cmd := exec.Command(m.notmuchPath, "tag", "+deleted", fmt.Sprintf("thread:%s", threadID))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete thread: %w", err)
	}
	return nil
}

// StarThread stars/unstars a thread
func (m *Manager) StarThread(threadID string, starred bool) error {
	tag := "+starred"
	if !starred {
		tag = "-starred"
	}

	cmd := exec.Command(m.notmuchPath, "tag", tag, fmt.Sprintf("thread:%s", threadID))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to star/unstar thread: %w", err)
	}
	return nil
}

// GetUnreadCount returns the total unread count
func (m *Manager) GetUnreadCount() (int, error) {
	cmd := exec.Command(m.notmuchPath, "count", "tag:unread")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	// Parse output to get count
	count := 0
	if _, err := fmt.Sscanf(strings.TrimSpace(string(output)), "%d", &count); err != nil {
		return 0, fmt.Errorf("failed to parse unread count: %w", err)
	}
	return count, nil
}

// GetMailFolders scans the mail directory and returns all available folders
func (m *Manager) GetMailFolders() ([]*MailFolder, error) {
	var folders []*MailFolder

	// Check if mail directory exists
	if _, err := os.Stat(m.maildirPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("mail directory does not exist: %s", m.maildirPath)
	}

	// Walk through the mail directory
	err := filepath.Walk(m.maildirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Log error but continue scanning other directories
			return nil
		}

		// Skip the root directory itself
		if path == m.maildirPath {
			return nil
		}

		// Only process directories
		if !info.IsDir() {
			return nil
		}

		// Get relative path from mail directory
		relPath, err := filepath.Rel(m.maildirPath, path)
		if err != nil {
			// Skip this directory if we can't get relative path
			return nil
		}

		// Skip hidden directories (starting with .)
		if strings.HasPrefix(relPath, ".") {
			return nil
		}

		// Skip Maildir storage folders (cur, new, tmp)
		if m.isMaildirStorageFolder(relPath) {
			return nil
		}

		// Check if this is a special folder
		isSpecial := m.isSpecialFolder(relPath)

		// Get unread and message counts using notmuch
		unreadCount, messageCount := m.getFolderCounts(relPath)

		folder := &MailFolder{
			Name:         relPath,
			Path:         path,
			UnreadCount:  unreadCount,
			MessageCount: messageCount,
			IsSpecial:    isSpecial,
		}

		folders = append(folders, folder)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan mail directory: %w", err)
	}

	// Sort folders: special folders first, then alphabetically
	sort.Slice(folders, func(i, j int) bool {
		if folders[i].IsSpecial && !folders[j].IsSpecial {
			return true
		}
		if !folders[i].IsSpecial && folders[j].IsSpecial {
			return false
		}
		return strings.ToLower(folders[i].Name) < strings.ToLower(folders[j].Name)
	})

	return folders, nil
}

// isSpecialFolder checks if a folder is a special system folder
func (m *Manager) isSpecialFolder(folderName string) bool {
	upperName := strings.ToUpper(folderName)
	specialFolders := []string{"INBOX", "SENT", "DRAFTS", "TRASH", "SPAM", "ARCHIVE", "JUNK"}

	for _, special := range specialFolders {
		if upperName == special {
			return true
		}
	}
	return false
}

// isMaildirStorageFolder checks if a folder is a Maildir storage folder
func (m *Manager) isMaildirStorageFolder(folderName string) bool {
	// Maildir storage folders that should not be displayed
	storageFolders := []string{"cur", "new", "tmp"}

	for _, storage := range storageFolders {
		if folderName == storage {
			return true
		}
	}
	return false
}

// getFolderCounts gets the unread and total message counts for a folder
func (m *Manager) getFolderCounts(folderName string) (unread, total int) {
	// Use notmuch to count messages in the folder
	// Format: folder:folderName
	query := fmt.Sprintf("folder:%s", folderName)

	// Get total count
	totalCmd := exec.Command(m.notmuchPath, "count", query)
	if output, err := totalCmd.Output(); err == nil {
		fmt.Sscanf(strings.TrimSpace(string(output)), "%d", &total)
	}

	// Get unread count
	unreadQuery := fmt.Sprintf("%s and tag:unread", query)
	unreadCmd := exec.Command(m.notmuchPath, "count", unreadQuery)
	if output, err := unreadCmd.Output(); err == nil {
		fmt.Sscanf(strings.TrimSpace(string(output)), "%d", &unread)
	}

	return unread, total
}

// parseNotmuchResults parses notmuch search results
// This is a placeholder - would need proper JSON parsing of notmuch output
func (m *Manager) parseNotmuchResults(output []byte) ([]*Thread, error) {
	// TODO: Implement proper parsing of notmuch JSON output
	// For now, return empty result
	return []*Thread{}, nil
}

// parseNotmuchThread parses notmuch thread output
// This is a placeholder - would need proper JSON parsing of notmuch output
func (m *Manager) parseNotmuchThread(output []byte) (*Thread, error) {
	// TODO: Implement proper parsing of notmuch thread output
	// For now, return empty thread
	return &Thread{}, nil
}
