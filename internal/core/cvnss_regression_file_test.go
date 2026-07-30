package core

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestCVNSSV120RegressionFile(t *testing.T) {
	f, err := os.Open("testdata/cvnss_regression_v120.tsv")
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
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) < 2 {
			t.Fatalf("invalid regression line: %q", line)
		}
		if got := DecodeCVNSS(parts[0]); got != parts[1] {
			t.Fatalf("DecodeCVNSS(%q)=%q want %q", parts[0], got, parts[1])
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count < 10 {
		t.Fatalf("expected at least 10 regression vectors, got %d", count)
	}
}
