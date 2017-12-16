//
// text.go implements a text structure and function for building
// an text editor functionality independent of UI.
//
package nuts

import (
	"fmt"
	"io"
)

const (
	// TextPlain describes something of type text/plain
	TextPlain = iota
	TextFDX
	TextMarkdown
	TextFountain
)

// Text is a structure for holding text content you
// want to process. The goal is to have a foundation
// for making stream processors, analyzers and editors
// easily.
//
// A Text structure should play nicely with any Reader, Writer
// and Seeker interfaces
type Text struct {
	// Name holds a human friendly name as a label
	Name string
	// Type is an integer value representing one of the text types supported (e.g. TextPlain).
	Type int
	// Metadata is a map to any additional data associated with the text, e.g. File Info, DOI content, etc.
	Metadata map[string]interface{}
	// Source holds the text in a byte slice
	Source []byte
	// Cursor holds the location of where to start reading or writing
	Cursor int64
}

//NOTE: Text holds content in memory. It implements io.Reader and io.Writer
//interfaces so it will be easily combined with the standard Go packages.

// Read copies t.Source into byte slice p up to the size of p
func (t *Text) Read(p []byte) (int, error) {
	sizeP := int64(len(p))
	sizeSrc := int64(len(t.Source))
	if t.Cursor == sizeSrc {
		return 0, io.EOF
	}
	n := 0
	//NOTE: we start reading from Cursor
	for i := t.Cursor; i < sizeP && i < sizeSrc; i++ {
		p[i] = t.Source[i]
		t.Cursor = i
		n++
	}
	return n, nil
}

// Write copies the content of p and appends it to t.Source
func (t *Text) Write(p []byte) (int, error) {
	sizeP := int64(len(p))
	sizeSrc := int64(len(t.Source))
	n := 0
	i := int64(0)
	// Overwrite our t.Source starting at t.Cursor until we need to
	// allocate a larger slice
	for ; i < sizeP && t.Cursor < sizeSrc; i++ {
		t.Source[t.Cursor] = p[i]
		t.Cursor++
		n++
	}
	// If sizeP puts us past the end of t.Source, append until we're done
	for ; i < sizeP; i++ {
		t.Source = append(t.Source, p[i])
		t.Cursor++
		n++
	}
	return n, nil
}

// Seek implements the io.Seeker interface for Text
func (t *Text) Seek(offset int64, whence int) (int64, error) {
	newCursor := int64(whence) + offset
	if newCursor < int64(len(t.Source)) {
		t.Cursor = newCursor
		return t.Cursor, nil
	}
	return t.Cursor, fmt.Errorf("offset %d whence %d is an invalid", offset, whence)
}

// ReadAt implements the io.ReadAt interface for Text
func (t *Text) ReadAt(p []byte, offset int64) (int, error) {
	_, err := t.Seek(offset, int(t.Cursor))
	if err != nil {
		return 0, err
	}
	return t.Read(p)
}

// WriteAt implements the io.WriteAt interface for Text
func (t *Text) WriteAt(p []byte, offset int64) (int, error) {
	_, err := t.Seek(offset, int(t.Cursor))
	if err != nil {
		return 0, err
	}
	return t.Write(p)
}
