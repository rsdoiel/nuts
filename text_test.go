package nuts

import (
	"bytes"
	"testing"
)

func TestReader(t *testing.T) {
	txt := new(Block)
	txt.Source = []byte("Hello World!!!!!")
	p := make([]byte, 5)
	n, err := txt.Read(p)
	if err != nil {
		t.Errorf("txt.Read(p): %s", err)
		t.FailNow()
	}
	if n != len(p) {
		t.Errorf("expected length %d, got %d for %s", n, len(p), p)
	}
}

func TestWriter(t *testing.T) {
	testPhrases := [][]byte{
		[]byte("Hello World!!!!!"),
		[]byte("This is a new World!!!!!!"),
		[]byte("Hope things get better!!!!!!"),
	}
	txt := new(Block)
	for i, p := range testPhrases {
		sizeSrc := len(txt.Source)
		sizeP := len(p)
		txt.Write(p)
		if len(txt.Source) != (sizeSrc + sizeP) {
			t.Errorf("expected (%d) len %d, got %d, for %q + %q", i, sizeSrc+sizeP, len(txt.Source), txt.Source, p)
		}
	}
}

func TestSeeker(t *testing.T) {
	txt := new(Block)
	txt.Source = []byte(`This is the way the world ends
This is the way the world ends
This is the way the world ends
Not with a bang but a whimper
`)
	offset, err := txt.Seek(int64(5), 0)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	if offset != txt.Cursor {
		t.Errorf("seek(5,0) expected %d, got %d", offset, txt.Cursor)
		t.FailNow()
	}
	expected := []byte("is")
	buf := make([]byte, 2)
	txt.Read(buf)
	if bytes.Compare(buf, expected) != 0 {
		t.Errorf("expected %q, got %q", expected, buf)
	}
}
