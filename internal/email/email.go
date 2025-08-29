package email

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

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
