package wordwrap

import (
	"bytes"
	"io"
	"unicode"
	"unicode/utf8"

	"github.com/corestoreio/csfw/util/bufferpool"
)

const nbsp = 0xA0

// String wraps the given string within lim width in characters.
//
// Wrapping is currently naive and only happens at white-space. A future
// version of the library will implement smarter wrapping. This means that
// pathological cases can dramatically reach past the limit, such as a very
// long word.
func String(s string, lim uint) string {
	// Initialize a buffer with a slightly larger size to account for breaks
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	Fstring(buf, s, lim)
	return buf.String()
}

var lineBreak = []byte("\n")

// Fstring same as String but writes into a buffer
func Fstring(buf io.Writer, s string, lim uint) {
	var current uint
	var wordBuf, spaceBuf bytes.Buffer

	var p [4]byte
	for _, char := range s {
		if char == '\n' {
			if wordBuf.Len() == 0 {
				if current+uint(spaceBuf.Len()) > lim {
					current = 0
				} else {
					current += uint(spaceBuf.Len())
					spaceBuf.WriteTo(buf)
				}
				spaceBuf.Reset()
			} else {
				current += uint(spaceBuf.Len() + wordBuf.Len())
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}

			ul := utf8.EncodeRune(p[:], char)
			buf.Write(p[:ul])
			current = 0
		} else if unicode.IsSpace(char) && char != nbsp {
			if spaceBuf.Len() == 0 || wordBuf.Len() > 0 {
				current += uint(spaceBuf.Len() + wordBuf.Len())
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}

			spaceBuf.WriteRune(char)
		} else {

			wordBuf.WriteRune(char)

			if current+uint(spaceBuf.Len()+wordBuf.Len()) > lim && uint(wordBuf.Len()) < lim {
				buf.Write(lineBreak)
				current = 0
				spaceBuf.Reset()
			}
		}
	}

	if wordBuf.Len() == 0 {
		if current+uint(spaceBuf.Len()) <= lim {
			spaceBuf.WriteTo(buf)
		}
	} else {
		spaceBuf.WriteTo(buf)
		wordBuf.WriteTo(buf)
	}
}
