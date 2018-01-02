package drumbeat

import (
	"encoding/gob"
	"io"
)

// WriteTo serializes the passed patterns and write them to writer.
func WriteTo(w io.Writer, patterns ...*Pattern) error {
	for _, p := range patterns {
		p.compact()
	}
	err := gob.NewEncoder(w).Encode(patterns)
	for _, p := range patterns {
		p.ReAlign()
	}
	return err
}

// ReadFrom reads a serialize drumbeat and returns the patterns
func ReadFrom(r io.Reader) ([]*Pattern, error) {
	var patterns []*Pattern
	err := gob.NewDecoder(r).Decode(&patterns)
	for _, p := range patterns {
		p.ReAlign()
	}
	return patterns, err
}
