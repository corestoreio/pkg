package pseudo

import (
	"bytes"
	"strings"
)

// TODO(Cyrill) refactor and user everywhere strings.Builder

// Character generates random character in the given language
func (s *Service) Character() string {
	return s.lookup(s.o.Lang, "characters", true)
}

// CharactersN generates n random characters in the given language
func (s *Service) CharactersN(n int) string {
	var chars []string
	for i := 0; i < n; i++ {
		chars = append(chars, s.Character())
	}
	return strings.Join(chars, "")
}

// Characters generates from 1 to 5 characters in the given language
func (s *Service) Characters() string {
	return s.CharactersN(s.r.Intn(5) + 1)
}

// Word generates random word
func (s *Service) Word(maxLen int) string {
	wrd := s.lookup(s.o.Lang, "words", true)
	if len(wrd) > int(maxLen) {
		wrd = wrd[:maxLen]
	}
	return wrd
}

// WordsN generates n random words
func (s *Service) WordsN(n int, maxLen int) string {
	words := make([]string, n)
	for i := 0; i < n; i++ {
		words[i] = s.Word(maxLen / n)
	}
	wrd := strings.Join(words, " ")
	if len(wrd) > int(maxLen) {

	}
	return wrd
}

// Words generates from 1 to 5 random words
func (s *Service) Words(maxLen int) string {
	return s.WordsN(s.r.Intn(5)+1, maxLen)
}

// Title generates from 2 to 5 titleized words
func (s *Service) Title() string {
	return strings.ToTitle(s.WordsN(2+s.r.Intn(4), 30))
}

// Sentence generates random sentence. If maxLen > 0, it will be considered as
// the total length of one sentence. The sentence gets cut after maxLen bytes.
func (s *Service) Sentence(maxLen int) string {
	var buf bytes.Buffer
	l := 3 + s.r.Intn(12)
	for i := 0; i < l; i++ {
		word := s.Word(maxLen / l)
		if s.r.Intn(5) == 0 && i+2 < l {
			word += ","
		}
		buf.WriteString(word)
		if i+1 < l {
			buf.WriteByte(' ')
		}
	}

	if s.r.Intn(8) == 0 {
		buf.WriteByte('!')
	} else {
		buf.WriteByte('.')
	}
	if ml := int(maxLen); ml > 0 && buf.Len() > ml {
		buf.Truncate(ml)
	}
	return buf.String()
}

// SentencesN generates n random sentences
func (s *Service) SentencesN(count int, maxLen int) string {
	sentences := make([]string, count)
	for i := 0; i < count; i++ {
		sentences[i] = s.Sentence(maxLen / count)
	}
	return strings.Join(sentences, " ")
}

// Sentences generates from 1 to 5 random sentences
func (s *Service) Sentences(maxLen int) string {
	return s.SentencesN(s.r.Intn(5)+1, maxLen)
}

// Paragraph generates paragraph
func (s *Service) Paragraph(maxLen int) string {
	return s.SentencesN(s.r.Intn(10)+1, maxLen)
}

// ParagraphsN generates n paragraphs
func (s *Service) ParagraphsN(n int, maxLen int) string {
	var paragraphs []string
	for i := 0; i < n; i++ {
		paragraphs = append(paragraphs, s.Paragraph(maxLen/n))
	}
	return strings.Join(paragraphs, "\t")
}

// Paragraphs generates from 1 to 5 paragraphs
func (s *Service) Paragraphs(maxLen int) string {
	return s.ParagraphsN(s.r.Intn(5)+1, maxLen)
}
