package gowal

type WalEntry struct {
	Idx   uint64
	Key   string
	Value []byte
}

func (we *WalEntry) Index() uint64 {
	return we.Idx
}
