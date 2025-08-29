package search

import (
	"fmt"
	"strings"
	"time"

	"github.com/romaintb/mel/internal/email"
)

// SearchType represents the type of search being performed
type SearchType int

const (
	SearchContent SearchType = iota
	SearchSender
	SearchGlobal
)

// SearchQuery represents a search query
type SearchQuery struct {
	Type      SearchType
	Query     string
	Filters   map[string]string
	SortBy    string
	SortOrder string
}

// SearchResult represents a search result with preview
type SearchResult struct {
	Thread    *email.Thread
	MatchType string
	MatchText string
	Context   string
	Relevance float64
}

// SearchService handles all search operations
type SearchService struct {
	emailManager *email.Manager
}

// NewSearchService creates a new search service
func NewSearchService(emailManager *email.Manager) *SearchService {
	return &SearchService{
		emailManager: emailManager,
	}
}

// Search performs a search based on the query type
func (s *SearchService) Search(query SearchQuery) ([]*SearchResult, error) {
	if s.emailManager == nil {
		return nil, fmt.Errorf("search service not initialized: email manager is nil")
	}

	switch query.Type {
	case SearchContent:
		return s.searchContent(query)
	case SearchSender:
		return s.searchSender(query)
	case SearchGlobal:
		return s.searchGlobal(query)
	default:
		return nil, fmt.Errorf("unknown search type: %v", query.Type)
	}
}

// searchContent performs full-text content search
func (s *SearchService) searchContent(query SearchQuery) ([]*SearchResult, error) {
	// Use notmuch for content search
	notmuchQuery := fmt.Sprintf("body:%s", query.Query)
	if query.Filters["folder"] != "" {
		notmuchQuery += fmt.Sprintf(" folder:%s", query.Filters["folder"])
	}
	if query.Filters["date"] != "" {
		notmuchQuery += fmt.Sprintf(" date:%s", query.Filters["date"])
	}

	// Perform the search
	results, err := s.emailManager.SearchEmails(notmuchQuery)
	if err != nil {
		return nil, fmt.Errorf("content search failed: %w", err)
	}

	// Convert to search results with context
	var searchResults []*SearchResult
	for _, thread := range results.Threads {
		result := &SearchResult{
			Thread:    thread,
			MatchType: "content",
			MatchText: query.Query,
			Context:   s.generateContext(thread, query.Query),
			Relevance: s.calculateRelevance(thread, query.Query),
		}
		searchResults = append(searchResults, result)
	}

	// Sort by relevance
	s.sortByRelevance(searchResults)
	return searchResults, nil
}

// searchSender performs sender-based search
func (s *SearchService) searchSender(query SearchQuery) ([]*SearchResult, error) {
	// Use notmuch for sender search
	notmuchQuery := fmt.Sprintf("from:%s", query.Query)
	if query.Filters["folder"] != "" {
		notmuchQuery += fmt.Sprintf(" folder:%s", query.Filters["folder"])
	}

	// Perform the search
	results, err := s.emailManager.SearchEmails(notmuchQuery)
	if err != nil {
		return nil, fmt.Errorf("sender search failed: %w", err)
	}

	// Convert to search results
	var searchResults []*SearchResult
	for _, thread := range results.Threads {
		result := &SearchResult{
			Thread:    thread,
			MatchType: "sender",
			MatchText: query.Query,
			Context:   s.generateSenderContext(thread),
			Relevance: s.calculateSenderRelevance(thread, query.Query),
		}
		searchResults = append(searchResults, result)
	}

	// Sort by relevance
	s.sortByRelevance(searchResults)
	return searchResults, nil
}

// searchGlobal performs global search across all fields
func (s *SearchService) searchGlobal(query SearchQuery) ([]*SearchResult, error) {
	// Use notmuch for global search
	notmuchQuery := query.Query
	if query.Filters["folder"] != "" {
		notmuchQuery += fmt.Sprintf(" folder:%s", query.Filters["folder"])
	}
	if query.Filters["date"] != "" {
		notmuchQuery += fmt.Sprintf(" date:%s", query.Filters["date"])
	}

	// Perform the search
	results, err := s.emailManager.SearchEmails(notmuchQuery)
	if err != nil {
		return nil, fmt.Errorf("global search failed: %w", err)
	}

	// Convert to search results
	var searchResults []*SearchResult
	for _, thread := range results.Threads {
		result := &SearchResult{
			Thread:    thread,
			MatchType: "global",
			MatchText: query.Query,
			Context:   s.generateGlobalContext(thread, query.Query),
			Relevance: s.calculateGlobalRelevance(thread, query.Query),
		}
		searchResults = append(searchResults, result)
	}

	// Sort by relevance
	s.sortByRelevance(searchResults)
	return searchResults, nil
}

// generateContext generates context for content search results
func (s *SearchService) generateContext(thread *email.Thread, query string) string {
	if thread.LatestMessage == nil {
		return "No content available"
	}

	body := thread.LatestMessage.Body
	if len(body) > 200 {
		body = body[:200] + "..."
	}

	// Highlight the search query in context
	highlighted := strings.ReplaceAll(body, query, fmt.Sprintf("**%s**", query))
	return highlighted
}

// generateSenderContext generates context for sender search results
func (s *SearchService) generateSenderContext(thread *email.Thread) string {
	if thread.LatestMessage == nil {
		return "No sender information"
	}

	return fmt.Sprintf("From: %s", thread.LatestMessage.From)
}

// generateGlobalContext generates context for global search results
func (s *SearchService) generateGlobalContext(thread *email.Thread, query string) string {
	if thread.LatestMessage == nil {
		return "No content available"
	}

	// Try to find the query in subject first
	if strings.Contains(strings.ToLower(thread.Subject), strings.ToLower(query)) {
		return fmt.Sprintf("Subject: %s", thread.Subject)
	}

	// Then in sender
	if strings.Contains(strings.ToLower(thread.LatestMessage.From), strings.ToLower(query)) {
		return fmt.Sprintf("From: %s", thread.LatestMessage.From)
	}

	// Finally in body
	return s.generateContext(thread, query)
}

// calculateRelevance calculates relevance score for content search
func (s *SearchService) calculateRelevance(thread *email.Thread, query string) float64 {
	relevance := 0.0

	// Boost for unread messages
	if thread.UnreadCount > 0 {
		relevance += 10.0
	}

	// Boost for recent messages
	daysSince := time.Since(thread.Timestamp).Hours() / 24
	if daysSince < 1 {
		relevance += 20.0
	} else if daysSince < 7 {
		relevance += 10.0
	} else if daysSince < 30 {
		relevance += 5.0
	}

	// Boost for starred messages
	if thread.LatestMessage != nil && thread.LatestMessage.Starred {
		relevance += 15.0
	}

	// Boost for message count (more active conversations)
	relevance += float64(thread.MessageCount) * 0.5

	return relevance
}

// calculateSenderRelevance calculates relevance score for sender search
func (s *SearchService) calculateSenderRelevance(thread *email.Thread, query string) float64 {
	relevance := s.calculateRelevance(thread, query)

	// Additional boost for exact sender matches
	if thread.LatestMessage != nil {
		if strings.EqualFold(thread.LatestMessage.From, query) {
			relevance += 25.0
		} else if strings.Contains(strings.ToLower(thread.LatestMessage.From), strings.ToLower(query)) {
			relevance += 15.0
		}
	}

	return relevance
}

// calculateGlobalRelevance calculates relevance score for global search
func (s *SearchService) calculateGlobalRelevance(thread *email.Thread, query string) float64 {
	relevance := s.calculateRelevance(thread, query)

	// Boost for subject matches
	if strings.Contains(strings.ToLower(thread.Subject), strings.ToLower(query)) {
		relevance += 20.0
	}

	// Boost for sender matches
	if thread.LatestMessage != nil {
		if strings.Contains(strings.ToLower(thread.LatestMessage.From), strings.ToLower(query)) {
			relevance += 15.0
		}
	}

	return relevance
}

// sortByRelevance sorts search results by relevance score
func (s *SearchService) sortByRelevance(results []*SearchResult) {
	// Simple bubble sort for now - could be optimized
	for i := 0; i < len(results)-1; i++ {
		for j := 0; j < len(results)-i-1; j++ {
			if results[j].Relevance < results[j+1].Relevance {
				results[j], results[j+1] = results[j+1], results[j]
			}
		}
	}
}
