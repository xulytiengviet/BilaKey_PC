package core

import "testing"

func TestCVNSSCriticalReverseCollisions(t *testing.T) {
	cases := map[string]string{
		"ses": "sẽ",
		"hed": "hề",
		"tod": "tồ",
		"tof": "tộ",
		"cos": "cõ",
	}
	for in, want := range cases {
		if got := DecodeCVNSS(in); got != want {
			t.Errorf("DecodeCVNSS(%q)=%q want %q", in, got, want)
		}
	}
}

func TestCVNSSCanonicalExamples(t *testing.T) {
	cases := map[string]string{
		"toiy": "tôi",
		"iwy":  "yêu",
		"vidf": "việt",
		"ily":  "yên",
		"idb":  "yết",
		"tizb": "tiếng",
	}
	for in, want := range cases {
		if got := DecodeCVNSS(in); got != want {
			t.Errorf("DecodeCVNSS(%q)=%q want %q", in, got, want)
		}
	}
}

func TestCVNSSQUGlideToneRegression(t *testing.T) {
	cases := map[string]string{
		"qyl": "quỳ",
		"qyz": "quỷ",
		"qys": "quỹ",
		"qyj": "quý",
		"qyr": "quỵ",
	}
	for in, want := range cases {
		if got := DecodeCVNSS(in); got != want {
			t.Errorf("DecodeCVNSS(%q)=%q want %q", in, got, want)
		}
		inspection := InspectCVNSS(in)
		if inspection.Selected != want {
			t.Errorf("InspectCVNSS(%q).Selected=%q want %q", in, inspection.Selected, want)
		}
	}
}

func TestCVNSSOnsetAwareCandidateRanking(t *testing.T) {
	inspection := InspectCVNSS("vidf")
	if inspection.Selected != "việt" {
		t.Fatalf("selected=%q want việt", inspection.Selected)
	}
	if len(inspection.Details) < 2 {
		t.Fatalf("expected raw candidate audit trail, got %+v", inspection)
	}
	if inspection.Details[0].Score <= inspection.Details[1].Score {
		t.Fatalf("candidate ranking is not deterministic: %+v", inspection.Details)
	}
	if len(inspection.Candidates) != 1 || inspection.Candidates[0] != "việt" {
		t.Fatalf("orthographic candidates=%v want [việt]", inspection.Candidates)
	}
}
