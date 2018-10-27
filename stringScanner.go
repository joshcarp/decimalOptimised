package decimal

import (
	"bytes"
	"io"
	"strings"
	"unicode"
)

type stringScanner struct {
	reader *strings.Reader
	wid    int
	widSet bool
}

func (s *stringScanner) ReadRune() (r rune, size int, err error) {
	return s.reader.ReadRune()
}

func (s *stringScanner) UnreadRune() error {
	return s.reader.UnreadRune()
}

func (s *stringScanner) SkipSpace() {
	for {
		ch, _, err := s.ReadRune()
		if err != nil {
			break
		}
		if !unicode.IsSpace(ch) {
			s.UnreadRune()
			break
		}
	}
}

func (s *stringScanner) Token(skipSpace bool, f func(rune) bool) (token []byte, err error) {
	if skipSpace {
		s.SkipSpace()
	}

	var buf bytes.Buffer
	for {
		r, _, err := s.ReadRune()
		if err != nil {
			logicCheck(err == io.EOF, "%v == io.EOF", err)
			break
		}
		// A dirty hack to recognise ∞, which UTF-8-encodes as [226, 136, 158]
		if !f(r) {
			err := s.UnreadRune()
			logicCheck(err == nil, "%v", err)
			break
		}
		buf.WriteRune(r)
	}
	return buf.Bytes(), nil
}

func (s *stringScanner) Width() (wid int, ok bool) {
	return s.wid, s.widSet
}

func (s *stringScanner) Read(buf []byte) (n int, err error) {
	return s.reader.Read(buf)
}

func isLetterOrInf(r rune) bool {
	return unicode.IsLetter(r) || r == '∞'
}
