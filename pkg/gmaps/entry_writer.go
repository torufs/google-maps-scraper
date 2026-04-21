package gmaps

// EntryWriter is the common interface implemented by all output writers.
type EntryWriter interface {
	// Write persists a single entry. Implementations must be safe to call
	// with a nil entry (they should silently no-op).
	Write(e *Entry) error

	// Close flushes any buffered data and releases resources.
	Close() error
}
