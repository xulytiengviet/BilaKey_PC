package core

import "testing"

func FuzzDecodeCVNSSDoesNotPanic(f *testing.F) {
	for _, seed := range []string{"toiy", "qyl", "vidf", "OpenAI", "", "điện", "😀"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		_ = DecodeCVNSS(input)
		_ = InspectCVNSS(input)
	})
}
