package core

import "testing"

func TestCVNSSDecode(t *testing.T) {
	cases := map[string]string{
		"toiy": "tôi",
		"iwy":  "yêu",
		"vidf": "việt",
		"ily":  "yên",
		"idb":  "yết",
	}
	for in, want := range cases {
		if got := DecodeCVNSS(in); got != want {
			t.Fatalf("DecodeCVNSS(%q)=%q want %q", in, got, want)
		}
	}
}

func TestCVNSSV500AuditInvariants(t *testing.T) {
	if CVNSSRuleVersion != "5.1.0-bilakey-core" {
		t.Fatalf("rule version = %q", CVNSSRuleVersion)
	}
	if CVNSSBaseRows != 758 || CVNSSPatchEntries != 336 || CVNSSAmbiguityGroups != 56 || CVNSSCanonicalPolicies != 56 || CVNSSCriticalCollisionGroups != 5 {
		t.Fatalf("unexpected oracle invariants: base=%d patches=%d ambiguity=%d policies=%d critical=%d", CVNSSBaseRows, CVNSSPatchEntries, CVNSSAmbiguityGroups, CVNSSCanonicalPolicies, CVNSSCriticalCollisionGroups)
	}
}

func TestCVNSSCandidateGraph(t *testing.T) {
	inspection := InspectCVNSS("ses")
	if !inspection.Ambiguous || !inspection.Critical || inspection.Selected != "sẽ" {
		t.Fatalf("unexpected inspection: %+v", inspection)
	}
	want := []string{"sẽ", "soec"}
	if len(inspection.Candidates) != len(want) {
		t.Fatalf("candidates=%v want=%v", inspection.Candidates, want)
	}
	for i := range want {
		if inspection.Candidates[i] != want[i] {
			t.Fatalf("candidates=%v want=%v", inspection.Candidates, want)
		}
	}
	copyOfCandidates := CVNSSCandidates("ses")
	copyOfCandidates[0] = "mutated"
	if InspectCVNSS("ses").Candidates[0] != "sẽ" {
		t.Fatal("candidate API leaked mutable generated data")
	}
}

func TestCVNSSUnambiguousCandidate(t *testing.T) {
	candidates := CVNSSCandidates("toiy")
	if len(candidates) != 1 || candidates[0] != "tôi" {
		t.Fatalf("candidates=%v want=[tôi]", candidates)
	}
}
