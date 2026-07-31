package hotkey

import "testing"

func TestResolve(t *testing.T) {
	cases := []struct {
		vk                    uint32
		ctrl, shift, alt, win bool
		want                  Action
	}{
		{0x20, true, true, false, false, ToggleVietnamese},
		{'1', true, true, false, false, SelectCVNSS},
		{'2', true, true, false, false, SelectVNITelex},
		{'3', true, true, false, false, None}, // 2.5.0 has only two input modes.
		{'0', true, true, false, false, CycleCandidate},
		{0x09, true, false, false, false, None}, // Ctrl+Tab belongs to the foreground app.
		{'1', true, true, true, false, None},
		{'1', true, true, false, true, None},
	}
	for _, tc := range cases {
		if got := Resolve(tc.vk, tc.ctrl, tc.shift, tc.alt, tc.win); got != tc.want {
			t.Errorf("Resolve(%#x)=%d want %d", tc.vk, got, tc.want)
		}
	}
}
