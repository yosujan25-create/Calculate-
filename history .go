package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// HistoryManager owns all FILE HANDLING for saving/loading history, and a
// sync.Mutex so it's safe to call from multiple goroutines at once.
type HistoryManager struct {
	mu       sync.Mutex
	FilePath string
}

func NewHistoryManager(path string) *HistoryManager {
	return &HistoryManager{FilePath: path}
}

// SaveHistory appends entries to a text file, creating the file/folder if needed.
func (h *HistoryManager) SaveHistory(entries []string) error {
	h.mu.Lock() // guard the file against concurrent writers
	defer h.mu.Unlock()

	if dir := filepath.Dir(h.FilePath); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("could not create data directory: %w", err)
		}
	}

	file, err := os.OpenFile(h.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open history file: %w", err)
	}
	defer file.Close() // always release the OS file handle

	writer := bufio.NewWriter(file)
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	for _, entry := range entries { // loop over entries, write one line each
		line := fmt.Sprintf("[%s] %s\n", timestamp, entry)
		if _, err := writer.WriteString(line); err != nil {
			return fmt.Errorf("write failed: %w", err)
		}
	}

	return writer.Flush() // buffered writer: must flush or data may be lost
}

// LoadHistory reads every saved line back out of the file.
func (h *HistoryManager) LoadHistory() ([]string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	file, err := os.Open(h.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil // no history yet is not an error condition
		}
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { // loop line by line
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
