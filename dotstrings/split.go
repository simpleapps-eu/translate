package dotstrings

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// IsFuzzyToken will return true for tokens that match the text "fuzzy".
// The match is case insensitive.
func IsFuzzyToken(token string) bool {
	return strings.EqualFold(token, "fuzzy")
}

// Split will split the file into (fuzzy, context, id, string) tuples.
//
// TODO count processed runes so we can point to a location when there is an error.
func Split() bufio.SplitFunc {

	type LexFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)

	var (
		lexer LexFunc

		lexContext LexFunc
		lexID      LexFunc
		lexString  LexFunc
	)

	var skipTo = func(data []byte, atEOF bool, sep string) (advance int, err error) {

		datalen := len(data)
		seplen := len(sep)

		var r rune
		var size int
		for p := 0; p < datalen; p += size {
			dataP := data[p:]

			// r may contain the value utf8.RuneError
			r, size = utf8.DecodeRune(dataP)
			if unicode.IsSpace(r) {
				continue
			}

			// request more data if we are not at EOF and lacking enough data to check for the sep
			if p+seplen > datalen {
				if !atEOF {
					return // request more data
				}
				advance = datalen
				err = fmt.Errorf("reached end of file while looking for %q", sep)
				return
			}

			// check for sep and advance to first location after the sep
			if strings.HasPrefix(string(dataP), sep) {
				advance = p + seplen
				return
			}

			advance = p
			err = fmt.Errorf("expected to find %q, found %q", sep, r)
			return
		}

		if !atEOF {
			return // request more data
		}
		advance = datalen
		err = fmt.Errorf("reached end of file while looking for %q", sep)
		return
	}

	var collectTo = func(data []byte, atEOF bool, sep string) (advance int, token []byte, err error) {

		seplen := len(sep)

		p := strings.Index(string(data), sep)

		if p == -1 {
			if !atEOF {
				return // request more data
			}

			err = fmt.Errorf("reached end of file while looking for %q", sep)
			return
		}

		advance = p + seplen
		token = data[:p]
		return
	}

	var collectString = func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		datalen := len(data)
		skipNextRune := false

		var r rune
		var size int
		for p := 0; p < datalen; p += size {
			// r may contain the value utf8.RuneError
			r, size = utf8.DecodeRune(data[p:])

			if skipNextRune {
				skipNextRune = false
				continue
			}

			if r == '\\' {
				skipNextRune = true
				continue
			}

			if r == '"' {
				advance = p + size
				token = data[:p]
				return
			}
		}

		if !atEOF {
			return // request more data
		}

		advance = datalen
		err = errors.New("reached end of file while reading a string")
		return
	}

	lexContext = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		var offset int

		advance, err = skipTo(data, atEOF, "/* ")
		if advance == 0 || err != nil {
			// hitting EOF while looking for a context is not an error
			if atEOF && err != nil && advance == len(data) {
				err = nil
			}
			return
		}
		offset += advance

		advance, token, err = collectTo(data[offset:], atEOF, " */")
		if advance == 0 || err != nil {
			return // request more data or report error
		}
		offset += advance

		advance = offset
		if IsFuzzyToken(string(token)) {
			return // Remain in lexContext when we are returning a Fuzzy token.
		}

		// Switch to ID lexer and return Context token.
		lexer = lexID
		return
	}

	lexID = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		var offset int

		advance, err = skipTo(data, atEOF, "\"")
		if advance == 0 || err != nil {
			return // request more data or report error
		}
		offset += advance

		advance, token, err = collectString(data[offset:], atEOF)
		if advance == 0 || err != nil {
			return // request more data or report error
		}
		offset += advance

		advance = offset
		lexer = lexString
		return
	}

	lexString = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		var offset int

		advance, err = skipTo(data, atEOF, "=")
		if advance == 0 || err != nil {
			return // request more data or report error
		}
		offset += advance

		advance, err = skipTo(data[offset:], atEOF, "\"")
		if advance == 0 || err != nil {
			return // request more data or report error
		}
		offset += advance

		advance, token, err = collectString(data[offset:], atEOF)
		if advance == 0 || err != nil {
			return // request more data or report error
		}
		offset += advance

		advance, err = skipTo(data[offset:], atEOF, ";")
		if advance == 0 || err != nil {
			if err == nil {
				token = nil // discard result to force loading of additonal data.
			}
			return // request more data or report error
		}
		offset += advance

		advance = offset
		lexer = lexContext
		return
	}

	lexer = lexContext

	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		return lexer(data, atEOF)
	}
}
