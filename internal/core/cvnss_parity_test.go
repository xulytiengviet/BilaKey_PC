package core

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

var cvnssIntentionalOracleDivergence = map[string]string{
	"es":  "ẽ",
	"tes": "tẽ",
	"ed":  "ề",
	"ted": "tề",
	"od":  "ồ",
	"tod": "tồ",
	"of":  "ộ",
	"tof": "tộ",
	"os":  "õ",
	"tos": "tõ",
}

func TestCVNSSGoldenParity(t *testing.T) {
	f, err := os.Open("testdata/cvnss_golden.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			t.Fatalf("invalid golden line: %q", line)
		}
		got := DecodeCVNSS(parts[0])
		want := parts[1]
		if patched, ok := cvnssIntentionalOracleDivergence[parts[0]]; ok {
			want = patched
		}
		if got != want {
			t.Fatalf("parity mismatch input=%q got=%q want=%q", parts[0], got, want)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count < 1500 {
		t.Fatalf("expected at least 1500 parity vectors, got %d", count)
	}
}
