package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

type HistoryCommand struct {
	filePath string
	mutex    sync.Mutex
}

func NewHistoryCommand() (*HistoryCommand, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config directory: %w", err)
	}

	dirPath := filepath.Join(configDir, "oc-ai")
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create history directory: %w", err)
	}

	return &HistoryCommand{
		filePath: filepath.Join(dirPath, "history.json"),
		mutex:    sync.Mutex{},
	}, nil
}

func (h *HistoryCommand) AddToHistory(command string) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Remove duplicate CLI tool name if present
	if strings.HasPrefix(command, activeTool+" ") {
		command = strings.TrimPrefix(command, activeTool+" ")
	}

	entries := h.loadHistory()
	entries = append(entries, HistoryEntry{
		Timestamp: time.Now(),
		Command:   command,
		Tool:      activeTool,
	})

	// Enforce history limit
	if cfg != nil && cfg.HistoryLimit > 0 && len(entries) > cfg.HistoryLimit {
		entries = entries[len(entries)-cfg.HistoryLimit:]
	}

	return h.saveHistory(entries)
}

func (h *HistoryCommand) loadHistory() []HistoryEntry {
	data, err := os.ReadFile(h.filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("Warning: Failed to read history file: %v\n", err)
		}
		return []HistoryEntry{}
	}

	var entries []HistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		fmt.Printf("Warning: Failed to parse history file: %v\n", err)
		// Backup corrupted file
		backupPath := h.filePath + ".bak"
		if err := os.Rename(h.filePath, backupPath); err != nil {
			fmt.Printf("Warning: Failed to backup corrupted history file: %v\n", err)
		}
		return []HistoryEntry{}
	}
	return entries
}

func (h *HistoryCommand) saveHistory(entries []HistoryEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history entries: %w", err)
	}

	// Write to temporary file first
	tempFile := h.filePath + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary history file: %w", err)
	}

	// Rename temporary file to actual file (atomic operation)
	if err := os.Rename(tempFile, h.filePath); err != nil {
		os.Remove(tempFile) // Clean up temp file if rename fails
		return fmt.Errorf("failed to save history file: %w", err)
	}
	return nil
}

type HistoryEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Command   string    `json:"command"`
	Tool      string    `json:"tool"`
}

func init() {
	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "Show command history",
		Run: func(cmd *cobra.Command, args []string) {
			hc, err := NewHistoryCommand()
			if err != nil {
				fmt.Printf("Error initializing history: %v\n", err)
				return
			}

			entries := hc.loadHistory()
			if len(entries) == 0 {
				fmt.Println("No command history found.")
				return
			}

			fmt.Println("Command History:")
			for i, entry := range entries {
				// Format command with tool name
				fmt.Printf("%d. [%s] %s %s\n",
					i+1,
					entry.Timestamp.Format("2006-01-02 15:04:05"),
					entry.Tool,
					entry.Command)
			}
		},
	}

	rootCmd.AddCommand(historyCmd)
}
