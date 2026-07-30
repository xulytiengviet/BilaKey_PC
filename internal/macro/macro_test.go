package macro

import (
	"path/filepath"
	"testing"
)

func TestTableRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "macros.tsv")
	tab := New(path)
	tab.Upsert("ko", "không")
	tab.Upsert("dc", "được")
	if err := tab.Save(); err != nil {
		t.Fatal(err)
	}
	loaded := New(path)
	if err := loaded.Load(); err != nil {
		t.Fatal(err)
	}
	if got, ok := loaded.Lookup("ko"); !ok || got != "không" {
		t.Fatalf("got %q ok=%v", got, ok)
	}
}
