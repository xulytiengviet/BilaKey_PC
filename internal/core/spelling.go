package core

import (
	"strings"
	"unicode"
)

var validInitials = []string{"ngh", "ng", "ch", "gh", "kh", "nh", "ph", "th", "tr", "gi", "qu", "b", "c", "d", "đ", "g", "h", "k", "l", "m", "n", "r", "s", "t", "v", "x"}
var validCodas = []string{"ch", "ng", "nh", "c", "m", "n", "p", "t", ""}

var validNuclei = map[string]bool{
	"a": true, "ă": true, "â": true, "e": true, "ê": true, "i": true, "o": true, "ô": true, "ơ": true, "u": true, "ư": true, "y": true,
	"ai": true, "ao": true, "au": true, "ay": true, "âu": true, "ây": true, "eo": true, "êu": true, "ia": true, "iê": true, "iêu": true,
	"iu": true, "oa": true, "oai": true, "oao": true, "oay": true, "oe": true, "oeo": true, "oi": true, "ôi": true, "ơi": true,
	"ua": true, "uâ": true, "uây": true, "uê": true, "ui": true, "uô": true, "uôi": true, "uy": true, "uya": true, "uyê": true, "uyu": true,
	"ưa": true, "ưi": true, "ươ": true, "ươi": true, "ươu": true, "ưu": true, "yê": true, "yêu": true, "oo": true, "uơ": true,
}

func IsLikelyVietnameseSyllable(word string) bool {
	word = strings.ToLower(stripToneString(word))
	if word == "" {
		return false
	}
	for _, r := range word {
		if !unicode.IsLetter(r) {
			return false
		}
	}

	initial := ""
	for _, p := range validInitials {
		if strings.HasPrefix(word, p) {
			initial = p
			break
		}
	}
	rest := strings.TrimPrefix(word, initial)
	coda := ""
	for _, p := range validCodas {
		if p != "" && strings.HasSuffix(rest, p) {
			coda = p
			break
		}
	}
	nucleus := strings.TrimSuffix(rest, coda)
	if !validNuclei[nucleus] {
		return false
	}

	// Common orthographic constraints.
	if strings.HasPrefix(word, "q") && initial != "qu" {
		return false
	}
	if initial == "k" && nucleus != "i" && !strings.HasPrefix(nucleus, "iê") && nucleus != "e" && nucleus != "ê" && !strings.HasPrefix(nucleus, "y") {
		return false
	}
	if initial == "c" && (strings.HasPrefix(nucleus, "i") || strings.HasPrefix(nucleus, "e") || strings.HasPrefix(nucleus, "ê") || strings.HasPrefix(nucleus, "y")) {
		return false
	}
	if initial == "gh" && !(strings.HasPrefix(nucleus, "i") || strings.HasPrefix(nucleus, "e") || strings.HasPrefix(nucleus, "ê")) {
		return false
	}
	if initial == "g" && (strings.HasPrefix(nucleus, "e") || strings.HasPrefix(nucleus, "ê")) {
		return false
	}
	if initial == "ngh" && !(strings.HasPrefix(nucleus, "i") || strings.HasPrefix(nucleus, "e") || strings.HasPrefix(nucleus, "ê")) {
		return false
	}
	return true
}
