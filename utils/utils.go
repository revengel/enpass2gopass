package utils

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// GetHashFromBytes -
func GetHashFromBytes(in []byte) string {
	h := sha256.New()
	h.Write(in)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// GetHash -
func GetHash(in string) string {
	return GetHashFromBytes([]byte(in))
}

// TruncStr -
func TruncStr(in string, maxLen int) string {
	if len(in) <= maxLen {
		return in
	}
	return strings.TrimSuffix(in[:maxLen], "_")
}

// Indent -
func Indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

// InList -
func InList[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
