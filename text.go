//
// text.go implements a text structure and function for building
// a text editor functionality independent of UI.
//
package nuts

import (
	"fmt"
	"io"
)

const (
	// Plain describes a text/plain document
	Plain = iota
	// FDX a Final Draft XML document
	FDX
	// Markdown document
	Markdown
	// Fountain document
	Fountain
)

// Formatter is an interface that can return a string reprentation
// and a byte array representation of one or more blocks
type Formatter interface {
	String() string
	Parse([]byte) (*Block, error)
	Render(io.Writer, *Block) error
}

type Document struct {
	// Name holds a human friendly name as a label
	Name string
	// Type is an integer value representing one of the text types supported (e.g. TextPlain).
	Type int
	// Meta is a map to any additional data associated with the text, e.g. File Info, DOI content, etc.
	Meta map[string]interface{}
	// Cursor holds the document level cursor (byte count into document)
	Cursor int
	// Source holds an ordered list of Block pointers
	Source []*Block
}

//FIXME: need a way to move from a document cursor to block number and block cursor

//FIXME: need something to hold an undo queue

// Block is a structure for holding text content you
// want to process. The goal is to have a foundation
// for making stream processors, analyzers and editors
// easily.
//
// A Block structure should play nicely with any Reader, Writer
// and Seeker interfaces
type Block struct {
	// Source holds the text in a byte slice
	Source []byte
	// Cursor holds the location of where to start reading or writing
	Cursor int
}

//NOTE: Block holds content in memory. It implements io.Reader and io.Writer
//interfaces so it will be easily combined with the standard Go packages.

// Read copies b.Source into byte slice p up to the size of p
func (b *Block) Read(p []byte) (int, error) {
	sizeP := len(p)
	sizeSrc := len(b.Source)
	if b.Cursor == sizeSrc {
		return 0, io.EOF
	}
	n := 0
	//NOTE: we start reading from Cursor
	for i := 0; i < sizeP && b.Cursor < sizeSrc; i++ {
		p[i] = b.Source[b.Cursor]
		b.Cursor++
		n++
	}
	return n, nil
}

// Write copies the content of p and appends it to b.Source
func (b *Block) Write(p []byte) (int, error) {
	sizeP := len(p)
	sizeSrc := len(b.Source)
	n := 0
	i := 0
	// Overwrite our b.Source starting at b.Cursor until we need to
	// allocate a larger slice
	for ; i < sizeP && b.Cursor < sizeSrc; i++ {
		b.Source[b.Cursor] = p[i]
		b.Cursor++
		n++
	}
	// If sizeP puts us past the end of b.Source, append until we're done
	for ; i < sizeP; i++ {
		b.Source = append(b.Source, p[i])
		b.Cursor++
		n++
	}
	return n, nil
}

// Seek implements the io.Seeker interface for Block
func (b *Block) Seek(offset int, whence int) (int, error) {
	newCursor := whence + offset
	if newCursor < len(b.Source) {
		b.Cursor = newCursor
		return b.Cursor, nil
	}
	return b.Cursor, fmt.Errorf("offset %d whence %d is an invalid", offset, whence)
}

// ReadAt implements the io.ReadAt interface for Block
func (b *Block) ReadAt(p []byte, offset int) (int, error) {
	_, err := b.Seek(offset, b.Cursor)
	if err != nil {
		return 0, err
	}
	return b.Read(p)
}

// WriteAt implements the io.WriteAt interface for Block
func (b *Block) WriteAt(p []byte, offset int) (int, error) {
	_, err := b.Seek(offset, b.Cursor)
	if err != nil {
		return 0, err
	}
	return b.Write(p)
}
