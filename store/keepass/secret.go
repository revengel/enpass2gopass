package keepass

import (
	"github.com/tobischo/gokeepasslib/v3"
	"github.com/tobischo/gokeepasslib/v3/wrappers"
)

// Secret -
type Secret struct {
	gokeepasslib.Entry
}

func (s Secret) setKey(k, v string, sensitivity bool) {
	s.Values = append(s.Values, gokeepasslib.ValueData{
		Key: k,
		Value: gokeepasslib.V{
			Content:   v,
			Protected: wrappers.NewBoolWrapper(sensitivity),
		},
	})
}

func (s Secret) setKeyOrAlt(k, altK, v string, sensitivity bool) {
	if t := s.GetContent(k); t == "" {
		s.setKey(k, v, sensitivity)
		return
	}
	s.setKey(altK, v, sensitivity)
}

// NewSecret -
func NewSecret() *Secret {
	var sec = gokeepasslib.NewEntry()
	return &Secret{sec}
}
