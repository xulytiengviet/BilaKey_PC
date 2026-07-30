package core

import "unicode"

func transformVNI(raw string, oldStyle bool) string {
	out := make([]rune, 0, len([]rune(raw)))
	pendingTone := toneNone

	for _, ch := range []rune(raw) {
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
		out = append(out, ch)
	}
	return applyToneToWord(string(out), pendingTone, oldStyle)
}

func vniTone(key rune) tone {
	switch key {
	case '1':
		return toneAcute
	case '2':
		return toneGrave
	case '3':
		return toneHook
	case '4':
		return toneTilde
	case '5':
		return toneDot
	default:
		return toneNone
	}
}

func applyLastEligible(out []rune, mapping map[rune]rune) bool {
	for i := len(out) - 1; i >= 0; i-- {
		base := unicode.ToLower(plainVowel(out[i]))
		if target, ok := mapping[base]; ok {
			out[i] = replaceBaseRune(out[i], target)
			return true
		}
	}
	return false
}

func applyVNI7(out []rune) bool {
	if len(out) >= 2 {
		a := unicode.ToLower(plainVowel(out[len(out)-2]))
		b := unicode.ToLower(plainVowel(out[len(out)-1]))
		if a == 'u' && b == 'o' {
			out[len(out)-2] = replaceBaseRune(out[len(out)-2], 'ư')
			out[len(out)-1] = replaceBaseRune(out[len(out)-1], 'ơ')
			return true
		}
	}
	return applyLastEligible(out, map[rune]rune{'u': 'ư', 'o': 'ơ'})
}

func undoLastVNIShape(out []rune, mapping map[rune]rune) bool {
	for i := len(out) - 1; i >= 0; i-- {
		if target, ok := mapping[unicode.ToLower(out[i])]; ok {
			out[i] = replaceBaseRune(out[i], target)
			return true
		}
	}
	return false
}

func undoVNI7(out []rune) bool {
	if len(out) >= 2 {
		a := unicode.ToLower(out[len(out)-2])
		b := unicode.ToLower(out[len(out)-1])
		if a == 'ư' && b == 'ơ' {
			out[len(out)-2] = replaceBaseRune(out[len(out)-2], 'u')
			out[len(out)-1] = replaceBaseRune(out[len(out)-1], 'o')
			return true
		}
	}
	return undoLastVNIShape(out, map[rune]rune{'ư': 'u', 'ơ': 'o'})
}
