package core

import "testing"

func TestTelex(t *testing.T) {
	e := New(MethodTelex, Options{})
	cases := map[string]string{
		"tieengs": "tiếng",
		"vieetj":  "việt",
		"ddoongf": "đồng",
		"hoaf":    "hòa",
		"thuys":   "thúy",
		"uow":     "ươ",
	}
	for in, want := range cases {
		if got := e.Transform(in); got != want {
			t.Fatalf("Telex %q: got %q want %q", in, got, want)
		}
	}
}

func TestTelexOldToneStyle(t *testing.T) {
	e := New(MethodTelex, Options{OldToneStyle: true})
	if got := e.Transform("hoaf"); got != "hoà" {
		t.Fatalf("got %q", got)
	}
}

func TestVNI(t *testing.T) {
	e := New(MethodVNI, Options{})
	cases := map[string]string{
		"tieng61": "tiếng",
		"viet65":  "việt",
		"d9ong62": "đồng",
		"hoa2":    "hòa",
	}
	for in, want := range cases {
		if got := e.Transform(in); got != want {
			t.Fatalf("VNI %q: got %q want %q", in, got, want)
		}
	}
}

func TestTelexEscapeSequences(t *testing.T) {
	e := New(MethodTelex, Options{})
	cases := map[string]string{
		"ass":  "as",
		"aaa":  "aa",
		"aww":  "aw",
		"uoww": "uow",
		"asz":  "a",
		"az":   "az",
	}
	for in, want := range cases {
		if got := e.Transform(in); got != want {
			t.Errorf("Telex escape %q: got %q want %q", in, got, want)
		}
	}
}

func TestVNIEscapeSequences(t *testing.T) {
	e := New(MethodVNI, Options{})
	cases := map[string]string{
		"a11": "a1",
		"a66": "a6",
		"o77": "o7",
		"a88": "a8",
		"d99": "d9",
	}
	for in, want := range cases {
		if got := e.Transform(in); got != want {
			t.Errorf("VNI escape %q: got %q want %q", in, got, want)
		}
	}
}

func TestSpellRestore(t *testing.T) {
	e := New(MethodTelex, Options{SpellCheck: true, AutoRestoreWrongKey: true})
	if got := e.Transform("bbas"); got != "bbas" {
		t.Fatalf("invalid Vietnamese word should restore raw, got %q", got)
	}
}

func TestSmartCapitalizationCompatibility(t *testing.T) {
	cases := []struct {
		method string
		in     string
		want   string
	}{
		{MethodCVNSS, "Toiy", "Tôi"},
		{MethodCVNSS, "TOIY", "TÔI"},
		{MethodTelex, "Tieengs", "Tiếng"},
		{MethodTelex, "TIEENGS", "TIẾNG"},
		{MethodVNI, "Tieng61", "Tiếng"},
		{MethodVNI, "TIENG61", "TIẾNG"},
	}
	for _, tc := range cases {
		e := New(tc.method, Options{})
		if got := e.Transform(tc.in); got != tc.want {
			t.Fatalf("%s %q: got %q want %q", tc.method, tc.in, got, tc.want)
		}
	}
}

func TestCVNSSMixedTextSafe(t *testing.T) {
	e := New(MethodCVNSS, Options{SpellCheck: true, AutoRestoreWrongKey: true})
	for _, literal := range []string{"OpenAI", "GitHub", "JavaScript", "https"} {
		if got := e.Transform(literal); got != literal {
			t.Errorf("mixed-text literal %q changed to %q", literal, got)
		}
	}
	if got := e.Transform("qyl"); got != "quỳ" {
		t.Fatalf("CVNSS transform qyl=%q want quỳ", got)
	}
}

func TestTransformTextMixedContent(t *testing.T) {
	e := New(MethodCVNSS, Options{SpellCheck: true, AutoRestoreWrongKey: true})
	input := "qyl tizb vidf · OpenAI/GitHub"
	want := "quỳ tiếng việt · OpenAI/GitHub"
	if got := e.TransformText(input); got != want {
		t.Fatalf("TransformText=%q want %q", got, want)
	}
}

func TestCVNSSAuditAPI(t *testing.T) {
	a := AuditCVNSS()
	if a.Codes != 810 || a.AmbiguityGroups != 56 || a.CanonicalPolicies != 56 || a.CriticalCollisions != 5 || a.SilentOverwrite != 0 {
		t.Fatalf("unexpected audit: %+v", a)
	}
}

func TestVNITelexUnifiedMode(t *testing.T) {
	e := New(MethodVNITelex, Options{})
	cases := map[string]string{
		"tieengs": "tiếng", // Telex
		"tieng61": "tiếng", // VNI
		"ddoongf": "đồng",  // Telex
		"d9ong62": "đồng",  // VNI
		"d9oongf": "đồng",  // VNI đ + Telex ô/dấu
		"ddong62": "đồng",  // Telex đ + VNI ô/dấu
		"vieet5":  "việt",  // Telex ê + VNI dấu nặng
		"2026":    "2026",
	}
	for in, want := range cases {
		if got := e.Transform(in); got != want {
			t.Errorf("VNI/Telex %q: got %q want %q", in, got, want)
		}
	}
}

func TestLegacyMethodsNormalizeToUnifiedMode(t *testing.T) {
	for _, method := range []string{"Telex", "VNI", "Telex/VNI", "VNI/Telex", "vni telex"} {
		if got := New(method, Options{}).Method; got != MethodVNITelex {
			t.Errorf("Normalize %q=%q want %q", method, got, MethodVNITelex)
		}
	}
}
