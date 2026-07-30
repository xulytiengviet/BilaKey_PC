package core

import "unicode"

type tone int

const (
	toneNone tone = iota
	toneGrave
	toneAcute
	toneHook
	toneTilde
	toneDot
)

type vowelFamily struct {
	plain rune
	forms [6]rune
}

var vowelFamilies = []vowelFamily{
	{plain: 'a', forms: [6]rune{'a', 'à', 'á', 'ả', 'ã', 'ạ'}},
	{plain: 'ă', forms: [6]rune{'ă', 'ằ', 'ắ', 'ẳ', 'ẵ', 'ặ'}},
	{plain: 'â', forms: [6]rune{'â', 'ầ', 'ấ', 'ẩ', 'ẫ', 'ậ'}},
	{plain: 'e', forms: [6]rune{'e', 'è', 'é', 'ẻ', 'ẽ', 'ẹ'}},
	{plain: 'ê', forms: [6]rune{'ê', 'ề', 'ế', 'ể', 'ễ', 'ệ'}},
	{plain: 'i', forms: [6]rune{'i', 'ì', 'í', 'ỉ', 'ĩ', 'ị'}},
	{plain: 'o', forms: [6]rune{'o', 'ò', 'ó', 'ỏ', 'õ', 'ọ'}},
	{plain: 'ô', forms: [6]rune{'ô', 'ồ', 'ố', 'ổ', 'ỗ', 'ộ'}},
	{plain: 'ơ', forms: [6]rune{'ơ', 'ờ', 'ớ', 'ở', 'ỡ', 'ợ'}},
	{plain: 'u', forms: [6]rune{'u', 'ù', 'ú', 'ủ', 'ũ', 'ụ'}},
	{plain: 'ư', forms: [6]rune{'ư', 'ừ', 'ứ', 'ử', 'ữ', 'ự'}},
	{plain: 'y', forms: [6]rune{'y', 'ỳ', 'ý', 'ỷ', 'ỹ', 'ỵ'}},
}

var runeToFamily map[rune]struct {
	base rune
	t    tone
}

func init() {
	runeToFamily = make(map[rune]struct {
		base rune
		t    tone
	})
	for _, f := range vowelFamilies {
		for i, r := range f.forms {
			runeToFamily[r] = struct {
				base rune
				t    tone
			}{f.plain, tone(i)}
			runeToFamily[unicode.ToUpper(r)] = struct {
				base rune
				t    tone
			}{unicode.ToUpper(f.plain), tone(i)}
		}
	}
}

func isVowel(r rune) bool {
	_, ok := runeToFamily[r]
	return ok
}

func plainVowel(r rune) rune {
	if x, ok := runeToFamily[r]; ok {
		return x.base
	}
	return r
}

func toneOf(r rune) tone {
	if x, ok := runeToFamily[r]; ok {
		return x.t
	}
	return toneNone
}

func withTone(r rune, t tone) rune {
	x, ok := runeToFamily[r]
	if !ok || t < toneNone || t > toneDot {
		return r
	}
	upper := unicode.IsUpper(r)
	base := unicode.ToLower(x.base)
	for _, f := range vowelFamilies {
		if f.plain == base {
			out := f.forms[t]
			if upper {
				out = unicode.ToUpper(out)
			}
			return out
		}
	}
	return r
}

func stripToneString(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		runes[i] = withTone(r, toneNone)
	}
	return string(runes)
}

func toneTarget(word string, oldStyle bool) int {
	r := []rune(word)
	vowelIdx := make([]int, 0, 3)
	for i, ch := range r {
		if isVowel(ch) {
			vowelIdx = append(vowelIdx, i)
		}
	}
	if len(vowelIdx) == 0 {
		return -1
	}
	if len(vowelIdx) > 1 && len(r) >= 2 {
		lo := []rune(string(r))
		for i := range lo {
			lo[i] = unicode.ToLower(lo[i])
		}
		if len(lo) >= 2 && lo[0] == 'q' && plainVowel(lo[1]) == 'u' && vowelIdx[0] == 1 {
			vowelIdx = vowelIdx[1:]
		}
		if len(vowelIdx) > 1 && len(lo) >= 2 && lo[0] == 'g' && plainVowel(lo[1]) == 'i' && vowelIdx[0] == 1 {
			vowelIdx = vowelIdx[1:]
		}
	}
	if len(vowelIdx) == 1 {
		return vowelIdx[0]
	}
	for i := len(vowelIdx) - 1; i >= 0; i-- {
		idx := vowelIdx[i]
		switch unicode.ToLower(plainVowel(r[idx])) {
		case 'ă', 'â', 'ê', 'ô', 'ơ':
			return idx
		}
	}
	for i := 0; i+1 < len(vowelIdx); i++ {
		a := unicode.ToLower(plainVowel(r[vowelIdx[i]]))
		b := unicode.ToLower(plainVowel(r[vowelIdx[i+1]]))
		if a == 'ư' && b == 'ơ' {
			return vowelIdx[i+1]
		}
	}
	lastV := vowelIdx[len(vowelIdx)-1]
	if lastV < len(r)-1 {
		return lastV
	}
	if len(vowelIdx) >= 2 {
		a := unicode.ToLower(plainVowel(r[vowelIdx[0]]))
		b := unicode.ToLower(plainVowel(r[vowelIdx[1]]))
		if (a == 'o' && (b == 'a' || b == 'e')) || (a == 'u' && b == 'y') {
			if oldStyle {
				return vowelIdx[1]
			}
			return vowelIdx[0]
		}
	}
	return vowelIdx[0]
}

func applyToneToWord(word string, t tone, oldStyle bool) string {
	r := []rune(word)
	for i, ch := range r {
		r[i] = withTone(ch, toneNone)
	}
	if t == toneNone {
		return string(r)
	}
	idx := toneTarget(string(r), oldStyle)
	if idx >= 0 {
		r[idx] = withTone(r[idx], t)
	}
	return string(r)
}

func replaceBaseRune(r rune, target rune) rune {
	upper := unicode.IsUpper(r)
	if upper {
		return unicode.ToUpper(target)
	}
	return target
}
