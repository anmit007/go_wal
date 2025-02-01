package gowal

import "os"

type Wal struct {
	logFile *os.File
	index   map[uint64]int
}
