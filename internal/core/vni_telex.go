package core

import "unicode"

// transformVNITelex accepts Telex and VNI control keys in one composition
// engine. Users may switch convention between words and may combine a Telex
// shape key with a VNI tone key (or the reverse) inside one Vietnamese word.
func transformVNITelex(raw string, oldStyle bool) string {
	out := make([]rune, 0, len([]rune(raw)))
	pendingTone := toneNone

	for _, ch := range []rune(raw) {
		// VNI numeric controls.
		switch ch {
		case '0', '1', '2', '3', '4', '5':
			if hasVowel(out) {
				nextTone := vniTone(ch)
				if ch != '0' && pendingTone == nextTone {
					pendingTone = toneNone
					out = append(out, ch)
					continue
				}
				pendingTone = nextTone
				continue
			}
		case '6':
			if applyLastEligible(out, map[rune]rune{'a': 'â', 'e': 'ê', 'o': 'ô'}) {
				continue
			}
			if undoLastVNIShape(out, map[rune]rune{'â': 'a', 'ê': 'e', 'ô': 'o'}) {
				out = append(out, ch)
				continue
			}
		case '7':
			if applyVNI7(out) {
				continue
			}
			if undoVNI7(out) {
				out = append(out, ch)
				continue
			}
		case '8':
			if applyLastEligible(out, map[rune]rune{'a': 'ă'}) {
				continue
			}
			if undoLastVNIShape(out, map[rune]rune{'ă': 'a'}) {
				out = append(out, ch)
				continue
			}
		case '9':
			if len(out) > 0 && unicode.ToLower(out[len(out)-1]) == 'd' {
				out[len(out)-1] = replaceBaseRune(out[len(out)-1], 'đ')
				continue
			}
			if len(out) > 0 && unicode.ToLower(out[len(out)-1]) == 'đ' {
				out[len(out)-1] = replaceBaseRune(out[len(out)-1], 'd')
				out = append(out, ch)
				continue
			}
		}

		// Telex alphabetic controls.
		lo := unicode.ToLower(ch)
		switch lo {
		case 's', 'f', 'r', 'x', 'j', 'z':
			if hasVowel(out) {
				nextTone := telexTone(lo)
				if lo == 'z' {
					if pendingTone != toneNone {
						pendingTone = toneNone
						continue
					}
					if stripTelexShapes(out) {
						continue
					}
				} else if pendingTone == nextTone {
					pendingTone = toneNone
					out = append(out, ch)
					continue
				} else {
					pendingTone = nextTone
					continue
				}
			}
		case 'a', 'e', 'o', 'd':
			if len(out) > 0 {
				last := out[len(out)-1]
				if unicode.ToLower(last) == lo {
					switch lo {
					case 'a':
						out[len(out)-1] = replaceBaseRune(last, 'â')
					case 'e':
						out[len(out)-1] = replaceBaseRune(last, 'ê')
					case 'o':
						out[len(out)-1] = replaceBaseRune(last, 'ô')
					case 'd':
						out[len(out)-1] = replaceBaseRune(last, 'đ')
					}
					continue
				}
				if base, ok := undoRepeatedTelexShape(last, lo); ok {
					out[len(out)-1] = base
					out = append(out, ch)
					continue
				}
			}
		case 'w':
			if applyTelexW(out) {
				continue
			}
			if undoTelexW(out) {
				out = append(out, ch)
				continue
			}
		}
		out = append(out, ch)
	}

	return applyToneToWord(string(out), pendingTone, oldStyle)
}
