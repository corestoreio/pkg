package pseudo

import (
	"strings"
)

// Gender generates random gender
func (s *Service) Gender() string {
	return s.lookup(s.o.Lang, "genders", true)
}

// GenderAbbrev returns first downcased letter of the random gender
func (s *Service) GenderAbbrev() string {
	g := s.Gender()
	if g != "" {
		return strings.ToLower(string(g[0]))
	}
	return ""
}

// Language generates random human language
func (s *Service) Language() string {
	return s.lookup(s.o.Lang, "languages", true)
}
