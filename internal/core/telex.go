package core

import "unicode"

func transformTelex(raw string, oldStyle bool) string {
	out := make([]rune, 0, len([]rune(raw)))
	pendingTone := toneNone

	for _, ch := range []rune(raw) {
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

func telexTone(key rune) tone {
	switch key {
	case 's':
		return toneAcute
	case 'f':
		return toneGrave
	case 'r':
		return toneHook
	case 'x':
		return toneTilde
	case 'j':
		return toneDot
	default:
		return toneNone
	}
}

func undoRepeatedTelexShape(last, key rune) (rune, bool) {
	targets := map[rune]rune{'a': 'â', 'e': 'ê', 'o': 'ô', 'd': 'đ'}
	if unicode.ToLower(last) != targets[key] {
		return last, false
	}
	base := key
	if unicode.IsUpper(last) {
		base = unicode.ToUpper(base)
	}
	return base, true
}

func undoTelexW(out []rune) bool {
	if len(out) >= 2 {
		a, b := unicode.ToLower(out[len(out)-2]), unicode.ToLower(out[len(out)-1])
		if a == 'ư' && b == 'ơ' {
			out[len(out)-2] = replaceBaseRune(out[len(out)-2], 'u')
			out[len(out)-1] = replaceBaseRune(out[len(out)-1], 'o')
			return true
		}
	}
	if len(out) == 0 {
		return false
	}
	last := out[len(out)-1]
	switch unicode.ToLower(last) {
	case 'ă':
		out[len(out)-1] = replaceBaseRune(last, 'a')
	case 'ơ':
		out[len(out)-1] = replaceBaseRune(last, 'o')
	case 'ư':
		out[len(out)-1] = replaceBaseRune(last, 'u')
	default:
		return false
	}
	return true
}

func stripTelexShapes(out []rune) bool {
	changed := false
	for i, ch := range out {
		var base rune
		switch unicode.ToLower(ch) {
		case 'ă', 'â':
			base = 'a'
		case 'ê':
			base = 'e'
		case 'ô', 'ơ':
			base = 'o'
		case 'ư':
			base = 'u'
		case 'đ':
			base = 'd'
		default:
			continue
		}
		out[i] = replaceBaseRune(ch, base)
		changed = true
	}
	return changed
}

func applyTelexW(out []rune) bool {
	if len(out) == 0 {
		return false
	}
	if len(out) >= 2 {
		a, b := out[len(out)-2], out[len(out)-1]
		if unicode.ToLower(plainVowel(a)) == 'u' && unicode.ToLower(plainVowel(b)) == 'o' {
			out[len(out)-2] = replaceBaseRune(a, 'ư')
			out[len(out)-1] = replaceBaseRune(b, 'ơ')
			return true
		}
	}
	last := out[len(out)-1]
	switch unicode.ToLower(plainVowel(last)) {
	case 'a':
		out[len(out)-1] = replaceBaseRune(last, 'ă')
		return true
	case 'o':
		out[len(out)-1] = replaceBaseRune(last, 'ơ')
		return true
	case 'u':
		out[len(out)-1] = replaceBaseRune(last, 'ư')
		return true
	}
	return false
}

func hasVowel(r []rune) bool {
	for _, ch := range r {
		if isVowel(ch) {
			return true
		}
	}
	return false
}
