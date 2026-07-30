package core

import (
	"sort"
	"strings"
	"unicode"
)

// CVNSSCandidate describes one auditable decoding path. The runtime keeps the
// complete candidate graph, but ranks candidates with onset-aware Vietnamese
// orthography instead of resolving the rime in isolation.
type CVNSSCandidate struct {
	Vowel        string
	Value        string
	Canonical    bool
	Orthographic bool
	Score        int
	Reason       string
}

// CVNSSInspection exposes the deterministic decision made for one code word.
// Candidates always contains at least the selected value. Details preserves
// the ranked audit trail used by the engine and the candidate UI.
type CVNSSInspection struct {
	Input      string
	Initial    string
	CodeVowel  string
	Ambiguous  bool
	Critical   bool
	Candidates []string
	Selected   string
	Details    []CVNSSCandidate
}

func DecodeCVNSS(word string) string { return decodeCodeWord(word, cvnssToCQN, true) }
func DecodeCVN(word string) string   { return decodeCodeWord(word, cvnToCQN, false) }

func decodeCodeWord(word string, vowelMap map[string]string, auditSafe bool) string {
	if word == "" {
		return ""
	}
	lower := strings.ToLower(word)
	initial, codeVowel := splitCodeWord(lower, vowelMap)
	cqnVowel := vowelMap[codeVowel]
	// Keep the hot path allocation-free for the 754 unambiguous codes. Only
	// enter the audit-aware ranker for a real candidate graph collision.
	if auditSafe {
		if _, ambiguous := cvnssCandidates[codeVowel]; ambiguous {
			if ranked := rankCVNSSCandidates(initial, codeVowel); len(ranked) > 0 {
				cqnVowel = ranked[0].Vowel
			}
		}
	}
	if cqnVowel == "" {
		cqnVowel = codeVowel
	}
	out := buildCQNWord(initial, cqnVowel)
	return applyWordCase(word, out)
}

// InspectCVNSS returns the orthographic candidates, the raw audit details and
// the onset-aware canonical selection.
func InspectCVNSS(word string) CVNSSInspection {
	if word == "" {
		return CVNSSInspection{}
	}
	lower := strings.ToLower(word)
	initial, codeVowel := splitCodeWord(lower, cvnssToCQN)
	details := rankCVNSSCandidates(initial, codeVowel)
	if len(details) == 0 {
		vowel := cvnssToCQN[codeVowel]
		if vowel == "" {
			vowel = codeVowel
		}
		value := buildCQNWord(initial, vowel)
		details = []CVNSSCandidate{{
			Vowel:        vowel,
			Value:        value,
			Canonical:    true,
			Orthographic: IsLikelyVietnameseSyllable(value),
			Score:        0,
			Reason:       "literal fallback",
		}}
	}

	words := make([]string, 0, len(details))
	for _, detail := range details {
		candidate := applyWordCase(word, detail.Value)
		if detail.Orthographic && !containsString(words, candidate) {
			words = append(words, candidate)
		}
	}
	// Never make the candidate API empty. This is important for literal input,
	// foreign words and forward-compatible rule extensions.
	if len(words) == 0 {
		for _, detail := range details {
			candidate := applyWordCase(word, detail.Value)
			if !containsString(words, candidate) {
				words = append(words, candidate)
			}
		}
	}
	selected := applyWordCase(word, details[0].Value)
	if !containsString(words, selected) {
		words = append([]string{selected}, words...)
	}
	return CVNSSInspection{
		Input:      word,
		Initial:    initial,
		CodeVowel:  codeVowel,
		Ambiguous:  len(words) > 1,
		Critical:   len(cvnssCriticalCandidates[codeVowel]) > 0,
		Candidates: words,
		Selected:   selected,
		Details:    details,
	}
}

func CVNSSCandidates(word string) []string {
	result := InspectCVNSS(word).Candidates
	return append([]string(nil), result...)
}

func rankCVNSSCandidates(initial, codeVowel string) []CVNSSCandidate {
	vowels := cvnssCandidates[codeVowel]
	canonical := cvnssToCQN[codeVowel]
	if len(vowels) == 0 {
		if canonical == "" {
			return nil
		}
		vowels = []string{canonical}
	}

	ranked := make([]CVNSSCandidate, 0, len(vowels))
	for index, vowel := range vowels {
		value := buildCQNWord(initial, vowel)
		orthographic, structuralReason := cvnssOrthographicCandidate(initial, vowel, value)
		score := 0
		reasons := make([]string, 0, 4)
		if vowel == canonical {
			score += 100
			reasons = append(reasons, "oracle canonical")
		}
		if orthographic {
			score += 200
			reasons = append(reasons, "Vietnamese orthography")
		} else {
			score -= 400
			reasons = append(reasons, structuralReason)
		}
		// A code initial q represents Quốc-ngữ "qu". If the candidate starts
		// with a tone-bearing u, the legacy cleanup would delete the tone. Give
		// the structurally equivalent u+tone-on-y form priority. cleanCQNAfterJoin
		// also migrates the tone defensively, so old dictionaries remain safe.
		if initial == "q" {
			r := []rune(vowel)
			if len(r) > 1 && unicode.ToLower(plainVowel(r[0])) == 'u' {
				if toneOf(r[0]) == toneNone {
					score += 300
					reasons = append(reasons, "qu glide preserves tone")
				} else {
					score -= 150
					reasons = append(reasons, "tone migrated from qu glide")
				}
			}
		}
		// Stable ordering is the final tie breaker, preserving oracle order.
		score -= index
		ranked = append(ranked, CVNSSCandidate{
			Vowel:        vowel,
			Value:        value,
			Canonical:    vowel == canonical,
			Orthographic: orthographic,
			Score:        score,
			Reason:       strings.Join(reasons, "; "),
		})
	}
	sort.SliceStable(ranked, func(i, j int) bool { return ranked[i].Score > ranked[j].Score })
	return ranked
}

func cvnssOrthographicCandidate(initial, vowel, value string) (bool, string) {
	if value == "" {
		return false, "empty candidate"
	}
	// yê/yêm/yên/yêng... are zero-onset spellings. After a consonant, the
	// standard spelling uses iê/iêm/iên/iêng... instead (tiếng, việt, nghiêng).
	if initial != "" {
		r := []rune(vowel)
		if len(r) >= 2 && unicode.ToLower(plainVowel(r[0])) == 'y' && unicode.ToLower(plainVowel(r[1])) == 'ê' {
			return false, "yê-family is not used after a consonant onset"
		}
	}
	if !IsLikelyVietnameseSyllable(value) {
		return false, "not a likely Vietnamese syllable"
	}
	return true, ""
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func splitCodeWord(word string, vowelMap map[string]string) (string, string) {
	if _, ok := vowelMap[word]; ok {
		return "", word
	}
	for _, initial := range codeInitials {
		if !strings.HasPrefix(word, initial) {
			continue
		}
		rest := strings.TrimPrefix(word, initial)
		if rest != "" {
			if _, ok := vowelMap[rest]; ok {
				return initial, rest
			}
		}
	}
	for _, initial := range codeInitials {
		if strings.HasPrefix(word, initial) {
			return initial, strings.TrimPrefix(word, initial)
		}
	}
	return "", word
}

func buildCQNWord(cvnInitial, cqnVowel string) string {
	return cleanCQNAfterJoin(restoreCQNInitial(cvnInitial, cqnVowel), cqnVowel)
}

func firstBaseVowel(s string) rune {
	for _, r := range s {
		return unicode.ToLower(plainVowel(r))
	}
	return 0
}

func restoreCQNInitial(cvnInitial, cqnVowel string) string {
	base := firstBaseVowel(cqnVowel)
	switch cvnInitial {
	case "q":
		return "qu"
	case "j":
		return "gi"
	case "z":
		return "d"
	case "d":
		return "đ"
	case "f":
		return "ph"
	case "k":
		return "kh"
	case "w":
		if strings.ContainsRune("iêe", base) {
			return "ngh"
		}
		return "ng"
	case "g":
		if strings.ContainsRune("iêe", base) {
			return "gh"
		}
		return "g"
	case "c":
		if strings.ContainsRune("iêe", base) {
			return "k"
		}
		return "c"
	default:
		return cvnInitial
	}
}

var iToY = map[rune]rune{'i': 'y', 'ì': 'ỳ', 'ỉ': 'ỷ', 'ĩ': 'ỹ', 'í': 'ý', 'ị': 'ỵ'}

func cleanCQNAfterJoin(initial, vowel string) string {
	runes := []rune(vowel)
	if initial == "" && len(runes) > 0 && !strings.HasPrefix(vowel, "ia") {
		if y, ok := iToY[runes[0]]; ok {
			runes[0] = y
			vowel = string(runes)
		}
	}
	if initial == "gi" && len(runes) > 0 && unicode.ToLower(plainVowel(runes[0])) == 'i' {
		if isGiShortVowel(vowel) {
			return "g" + vowel
		}
		if len(runes) > 1 {
			vowel = string(runes[1:])
		} else {
			vowel = ""
		}
	}
	runes = []rune(vowel)
	if initial == "qu" && len(runes) > 1 && unicode.ToLower(plainVowel(runes[0])) == 'u' {
		// The u in qu is a glide and is removed when joining. Legacy candidate
		// forms such as ùy placed the tone on that glide; migrate it to the next
		// vowel before removal so qyl always becomes quỳ, never quy.
		if t := toneOf(runes[0]); t != toneNone {
			runes[1] = withTone(runes[1], t)
		}
		vowel = string(runes[1:])
	}
	return initial + vowel
}

func isGiShortVowel(v string) bool {
	r := []rune(v)
	if len(r) == 1 {
		return unicode.ToLower(plainVowel(r[0])) == 'i'
	}
	if len(r) == 2 && unicode.ToLower(plainVowel(r[0])) == 'i' {
		switch unicode.ToLower(r[1]) {
		case 'm', 'n', 'p', 't', 'c':
			return true
		}
	}
	if len(r) == 3 && unicode.ToLower(plainVowel(r[0])) == 'i' {
		tail := string([]rune{unicode.ToLower(r[1]), unicode.ToLower(r[2])})
		return tail == "ch" || tail == "ng" || tail == "nh"
	}
	return false
}

func applyWordCase(source, value string) string {
	if value == "" {
		return value
	}
	if source == strings.ToUpper(source) && source != strings.ToLower(source) {
		return strings.ToUpper(value)
	}
	sr := []rune(source)
	vr := []rune(value)
	if len(sr) > 0 && len(vr) > 0 && unicode.IsUpper(sr[0]) {
		vr[0] = unicode.ToUpper(vr[0])
		return string(vr)
	}
	return value
}
