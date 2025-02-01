package gowal

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

type Wal struct {
	LogFile *os.File
	Index   map[uint64]int64
}

func (w *Wal) OpenLogFile(filePath string) error {
	var err error
	w.LogFile, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 064)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	return nil
}
func (w *Wal) CloseLogFile() error {
	if w.LogFile != nil {
		if err := w.LogFile.Close(); err != nil {
			return fmt.Errorf("failed to close log file: %w", err)
		}
	}
	return nil
}

func (w *Wal) Write(index uint64, key string, value []byte) error {
	if w.LogFile == nil {
		return fmt.Errorf("LOG FILE IS NOT OPEN")
	}
	currentOffset, err := w.LogFile.Seek(0, os.SEEK_CUR)
	if err != nil {
		return fmt.Errorf("failed to get current offsets%w", err)
	}
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	entry := WalEntry{Idx: index, Key: key, Value: value}
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("failed to encode entry:%w", err)
	}
	if _, err := w.LogFile.Write(buff.Bytes()); err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}
	w.Index[index] = currentOffset
	return nil
}

// return key value and probably an error
func (w *Wal) Read(index uint64) (string, []byte, error) {
	offset, exists := w.Index[index]
	if !exists {
		return "", nil, fmt.Errorf("index not found")
	}
	if _, err := w.LogFile.Seek(offset, os.SEEK_SET); err != nil {
		return "", nil, fmt.Errorf("Failed to seek")
	}
	var entry WalEntry
	decode := gob.NewDecoder(w.LogFile)
	if err := decode.Decode(&entry); err != nil {
		return "", nil, fmt.Errorf("failed to decode entry %w", err)
	}
	return entry.Key, entry.Value, nil

}
func NewWAL(filePath string) (*Wal, error) {
	wal := &Wal{
		Index: make(map[uint64]int64),
	}

	// Open the log file
	var err error
	wal.LogFile, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return wal, nil
}
