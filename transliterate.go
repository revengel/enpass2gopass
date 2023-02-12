package main

import (
	"bytes"
	"regexp"
	"strings"
)

var baseRuEn = map[string]string{
	"а": "a",
	"б": "b",
	"в": "v",
	"г": "g",
	"д": "d",

	"е": "e",
	"ё": "e",
	"ж": "zh",
	"з": "z",
	"и": "i",

	"й": "о",
	"к": "k",
	"л": "l",
	"м": "m",
	"н": "n",

	"о": "o",
	"п": "p",
	"р": "r",
	"с": "s",
	"т": "t",

	"у": "u",
	"ф": "f",
	"х": "h",
	"ц": "c",
	"ч": "ch",

	"ш": "sh",
	"щ": "shch",
	"ъ": "",
	"ы": "y",
	"ь": "",

	"э": "e",
	"ю": "ju",
	"я": "ya",
}

func isRussianChar(r rune) bool {
	switch {
	case r >= 1040 && r <= 1103,
		r == 1105, r == 1025:
		return true
	}

	return false
}

func transliterate(in string) string {
	if in == "" {
		return ""
	}
	in = strings.ToLower(in)

	var input = bytes.NewBufferString(in)
	var output = bytes.NewBuffer(nil)

	var rr string
	var ok bool
	for {
		r, _, err := input.ReadRune()
		if err != nil {
			break
		}

		if !isRussianChar(r) {
			output.WriteRune(r)
			continue
		}

		rr, ok = baseRuEn[string(r)]
		if ok {
			output.WriteString(rr)
			continue
		}
	}

	in = output.String()
	var re = regexp.MustCompile(`[\W_]+`)
	in = re.ReplaceAllString(in, "_")
	return strings.Trim(in, "_")
}
