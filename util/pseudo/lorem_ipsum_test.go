package pseudo

import (
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestLoremIpsum(t *testing.T) {
	s := NewService(0, nil)
	for _, lang := range s.GetLangs() {
		assert.NoError(t, s.SetLang(lang))

		v := s.Sentence(120)
		assert.LenBetween(t, v, 1, 120, "Sentence failed with lang %s", lang)
		t.Logf("Sentence %q", v)

		v = s.Paragraph(410)
		assert.LenBetween(t, v, 1, 410, "Paragraph failed with lang %s", lang)
		t.Logf("Paragraph %q", v)

		v = s.ParagraphsN(2, 410)
		assert.LenBetween(t, v, 1, 420, "ParagraphsN failed with lang %s", lang)
		t.Logf("ParagraphsN %q", v)

		v = s.Paragraphs(1555)
		assert.LenBetween(t, v, 1, 1700, "Paragraphs failed with lang %s", lang)
		t.Logf("Paragraphs %q", v)

		v = s.Character()
		if v == "" {
			t.Errorf("Character failed with lang %s", lang)
		}

		v = s.CharactersN(2)
		if v == "" {
			t.Errorf("CharactersN failed with lang %s", lang)
		}

		v = s.Characters()
		if v == "" {
			t.Errorf("Characters failed with lang %s", lang)
		}

		v = s.Word(11)
		assert.LenBetween(t, v, 1, 565, "Paragraphs failed with lang %s", lang)
		t.Logf("Paragraphs %q", v)
		if v == "" {
			t.Errorf("Word failed with lang %s", lang)
		}

		v = s.WordsN(2, 40)
		assert.LenBetween(t, v, 1, 42, "WordsN failed with lang %s", lang)
		t.Logf("WordsN %q", v)

		v = s.Words(50)
		assert.LenBetween(t, v, 1, 55, "Words failed with lang %s", lang)
		t.Logf("Words %q", v)

		v = s.Title()
		if v == "" {
			t.Errorf("Title failed with lang %s", lang)
		}

		v = s.SentencesN(2, 200)
		assert.LenBetween(t, v, 1, 210, "SentencesN failed with lang %s", lang)
		t.Logf("SentencesN %q", v)

		v = s.Sentences(222)
		assert.LenBetween(t, v, 1, 240, "Sentences failed with lang %s", lang)
		t.Logf("Sentences %q", v)

	}
}
