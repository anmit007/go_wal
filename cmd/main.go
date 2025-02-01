package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	gowal "github.com/anmit007/go_wal"
)

func main() {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	// Create the log file path in home directory
	walPath := filepath.Join(homeDir, "wal.log")

	wal, err := gowal.NewWAL(walPath)
	if err != nil {
		log.Fatalf("Failed to initialize WAL: %v", err)
	}
	defer wal.CloseLogFile()
	entries := []struct {
		index uint64
		key   string
		value string
	}{
		{1, "key1", "value1"},
		{2, "key2", "value2"},
		{3, "key3", "value3"},
	}
	for _, entry := range entries {
		if err := wal.Write(entry.index, entry.key, []byte(entry.value)); err != nil {
			log.Fatalf("Failed to write entry %d: %v", entry.index, err)
		}
		fmt.Printf("Successfully wrote entry %d\n", entry.index)
	}

	// Read entries back
	for _, entry := range entries {
		key, value, err := wal.Read(entry.index)
		if err != nil {
			log.Fatalf("Failed to read entry %d: %v", entry.index, err)
		}
		fmt.Printf("Read entry %d: key=%s, value=%s\n", entry.index, key, string(value))
	}
}
