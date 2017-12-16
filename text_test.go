package nuts

import (
	"testing"
)

func TestReader(t *testing.T) {
	txt := new(Text)
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
	txt := new(Text)
	for i, p := range testPhrases {
		sizeSrc := len(txt.Source)
		sizeP := len(p)
		txt.Write(p)
		if len(txt.Source) != (sizeSrc + sizeP) {
			t.Errorf("expected (%d) len %d, got %d, for %q + %q", i, sizeSrc+sizeP, len(txt.Source), txt.Source, p)
		}
	}
}
