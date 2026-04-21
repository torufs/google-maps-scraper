package gmaps

// Writer is the common interface for all output format writers.
type Writer interface {
	// Write serialises a single Entry to the underlying output.
	// Implementations must treat a nil Entry as a no-op.
	Write(e *Entry) error
}

// Flusher is implemented by writers that buffer output and need an
// explicit flush step before the data is fully written (e.g. Excel).
type Flusher interface {
	Flush() error
}

// WriteAll writes all entries to w and, if w also implements Flusher,
// flushes the buffer afterwards.
func WriteAll(w Writer, entries []*Entry) error {
	for _, e := range entries {
		if err := w.Write(e); err != nil {
			return err
		}
	}
	if f, ok := w.(Flusher); ok {
		return f.Flush()
	}
	return nil
}
