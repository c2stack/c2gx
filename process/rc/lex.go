package rc

import (
	"errors"
	"strings"
	"unicode/utf8"
	"fmt"
	"github.com/c2gx/process"
)

// This uses the go feature call go tools in the build process. To ensure this gets
//  called before compilation, make this call before building
//
//    go generate c2gx/process/rc
//
//go:generate go tool yacc -o parser.go syntax.y
type Token struct {
	typ int
	val string
}

const eof rune = 0

const (
	ParseEof = iota + 1
	ParseErr
)

// KEEP LIST IN SYNC WITH syntax.y
var keywords = [...]string{
	"\n",
	";",
	"(",
	")",
	"{",
	"}",
	"&",
	"=",
	"^",
	"!",
	"[word]",
	"for",
	"in",
	"while",
	"if",
	"not",
	"~",
	"@",
	"switch",
	"fn",
}

type stateFunc func(*lexer) stateFunc

type lexer struct {
	pos       int
	start     int
	width     int
	state     stateFunc
	tree      process.Op
	input     string
	tokens    []Token
	head      int
	tail      int
	lastError error
}

func (l *lexer) acceptWS() {
	for {
		if l.next() != ' ' {
			l.backup()
			return
		}
		l.ignore()
		if l.isEof() {
			return
		}
	}
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) isEof() bool {
	return l.pos >= len(l.input)
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) popToken() Token {
	token := l.tokens[l.tail]
	l.tail = (l.tail + 1) % len(l.tokens)
	return token
}

func (l *lexer) pushToken(t Token) {
	l.tokens[l.head] = t
	l.head = (l.head + 1) % len(l.tokens)
}

func (l *lexer) nextToken() (Token, error) {
	for {
		if l.head != l.tail {
			token := l.popToken()
			if token.typ == ParseEof {
				return token, errors.New(token.val)
			}
			return token, nil
		} else {
			if l.state == nil {
				return Token{ParseEof, "EOF"}, nil
			}
			l.state = l.state(l)
		}
	}
}

func (l *lexer) emit(t int) {
	s := l.input[l.start:l.pos]
	l.pushToken(Token{t, s})
	l.start = l.pos

	// nothing to do w/emit but convienent place
	l.acceptWS()
}

func (l *lexer) keyword(ttype int) string {
	switch true {
	case ttype == ParseEof:
		return "EOF"
	case ttype == ParseErr:
		return "ERR"
	case ttype < token_eol:
		panic(fmt.Sprintf("%d is not a keyword", ttype))
	}
	return keywords[ttype-token_eol]
}

const nonWordRunes = "\n \t#;&|^$=`'{}()<>"

func (l *lexer) isWordRune(r rune) bool {
	return ! strings.ContainsRune(nonWordRunes, r)
}

func (l *lexer) acceptAlphaNumeric(typ int) bool {
	s := l.input[l.pos:]
	length := len(s)
	if length == 0 {
		return false
	}
	length = strings.IndexFunc(s, notRune(l.isWordRune))
	if length == -1 {
		length = len(s)
	}
	l.pos += length
	l.emit(typ)
	return true
}

func notRune(f func(rune) bool) func(rune) bool {
	return func(r rune) bool {
		return !f(r)
	}
}

func (l *lexer) acceptToken(typ int) bool {
	var keyword = l.keyword(typ)
	switch typ {
	case token_word:
		r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
		if l.isWordRune(r) {
			return l.acceptAlphaNumeric(typ)
		}
		return false
	}
	if !strings.HasPrefix(l.input[l.pos:], keyword) {
		return false
	}
	l.pos += len(keyword)
	l.emit(typ)
	return true
}

var looseTokens = [...]int {
	token_eol,
	token_semi,
	token_amp,
	token_eq,
	token_caret,
	token_bang,
}

var kywdTokens = []int {
	kywd_for,
	kywd_in,
	kywd_while,
	kywd_if,
	kywd_not,
	kywd_twiddle,
	kywd_subshell,
	kywd_switch,
	kywd_fn,
}

func (l *lexer) error(msg string) stateFunc {
	l.pushToken(Token{
		ParseErr,
		msg,
	})
	l.Error(msg)
	return nil
}

func (l *lexer) acceptWord() bool {
	// TODO: expand
	return l.acceptToken(token_word)
}

func beginLex(l *lexer) stateFunc {
	if l.acceptToken(kywd_for) {
		if ! l.acceptToken(token_open_paren) {
			return l.error("Expected '(' after 'for' statement")
		}
		if ! l.acceptWord() {
			return l.error("Expected identifier")
		}

		if l.acceptToken(kywd_in) {
			for l.acceptWord() {}
		}

		if ! l.acceptToken(token_closed_paren) {
			return l.error("Expected ')' in 'for' statement")
		}

		return beginLex
	}
	for _, tok := range kywdTokens {
		if l.acceptToken(tok) {
			return beginLex
		}
	}
	for _, tok := range looseTokens {
		if l.acceptToken(tok) {
			return beginLex
		}
	}
	if l.acceptToken(token_word) {
		return beginLex
	}
	return nil
}

func lex(input string) *lexer {
	l := &lexer{
		input:  input,
		tokens: make([]Token, 128),
		state:  beginLex,
	}
	l.acceptWS()
	l.state = beginLex(l)
	return l
}

func (l *lexer) Position() (line, col int) {
	for p := 0; p < l.pos; p++ {
		if l.input[p] == '\n' {
			line += 1
			col = 0
		} else {
			col += 1
		}
	}
	return
}

func (l *lexer) Error(e string) {
	line, col := l.Position()
	msg := fmt.Sprintf("%s - line %d, col %d", e, line, col)
	l.lastError = errors.New(msg)
}