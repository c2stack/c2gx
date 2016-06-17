package rc

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// This uses the go feature call go tools in the build process. To ensure this gets
//  called before compilation, make this call before building
//
//    go generate c2gx/process/rc
//
//go:generate go tool yacc -o parser.go rc.y
type Token struct {
	typ int
	val string
}

const eof rune = 0

const (
	char_doublequote = '"'
	char_singlequote = '\''
	char_backslash   = '\\'
)

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
	"[ident]",
	"[string]",
	"[num]",
	"for",
	"in",
	"while",
	"if",
	"not",
	"~",
	"@",
	"switch",
	"fn",
	"`",
}

type stateFunc func(*lexer) stateFunc

type lexer struct {
	pos       int
	start     int
	width     int
	state     stateFunc
	stack     *opStack
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

// Any of String, number or ident
func (l *lexer) acceptValue() bool {
	if !l.acceptString(token_string) {
		if !l.acceptNumber(token_number) {
			return l.acceptIdent(token_ident)
		}
	}
	return false
}

func (l *lexer) acceptNumber(typ int) bool {
	for i := 0; true; i++ {
		r := l.next()
		if unicode.IsNumber(r) {
			continue
		} else if i == 0 {
			break
		}
		l.backup()
		l.emit(typ)
		return true
	}
	l.pos = l.start
	return false

}

func (l *lexer) acceptIdent(typ int) bool {
	for i := 0; true; i++ {
		r := l.next()
		if (unicode.IsDigit(r) && i > 0) ||
			unicode.IsLetter(r) ||
			(r == '-' && i > 0) ||
			(r == '.' && i > 0) ||
			r == '_' {
			continue
		}
		if i == 0 {
			break
		}
		l.backup()
		l.emit(typ)
		return true
	}
	l.pos = l.start
	return false
}

func (l *lexer) acceptString(typ int) bool {
	r := l.next()
	if r != char_singlequote {
		l.backup()
		return false
	}
	for {
		r = l.next()
		if r == char_backslash {
			l.next()
		} else if r == char_singlequote {
			l.emit(typ)
			return true
		} else if r == eof {
			// bad format?
			return false
		}
	}
}

func notRune(f func(rune) bool) func(rune) bool {
	return func(r rune) bool {
		return !f(r)
	}
}

func (l *lexer) acceptToken(typ int) bool {
	var keyword = l.keyword(typ)
	switch typ {
	case token_ident:
		return l.acceptIdent(typ)
	case token_string:
		return l.acceptString(typ)
	case token_number:
		return l.acceptNumber(typ)
	}
	if !strings.HasPrefix(l.input[l.pos:], keyword) {
		return false
	}
	l.pos += len(keyword)
	l.emit(typ)
	return true
}

var looseTokens = [...]int{
	token_eol,
	token_semi,
	token_amp,
	token_eq,
	token_caret,
	token_bang,
	token_in,
	token_while,
	token_if,
	token_not,
	token_twiddle,
	token_subshell,
	token_switch,
	token_number,
	token_ident,
	token_string,
	token_closed_brace,
	token_backtick,
}

func (l *lexer) error(msg string) stateFunc {
	l.pushToken(Token{
		ParseErr,
		msg,
	})
	l.Error(msg)
	return nil
}

func beginLex(l *lexer) stateFunc {
	if l.acceptToken(token_for) {
		if !l.acceptToken(token_open_paren) {
			return l.error("Expected '(' after 'for' statement")
		}
		if !l.acceptIdent(token_ident) {
			return l.error("Expected identifier")
		}

		if l.acceptToken(token_in) {
			for {
				if !l.acceptValue() {
					if !l.acceptToken(token_closed_paren) {
						return l.error("Expected value or ')'")
					}
				}
			}
		}

		if !l.acceptToken(token_closed_paren) {
			return l.error("Expected ')' in 'for' statement")
		}

		return beginLex
	}
	if l.acceptToken(token_fn) {
		if !l.acceptIdent(token_ident) {
			return l.error("Expected identifier after fn")
		}
		if !l.acceptToken(token_open_brace) {
			return l.error("Expected '{' after fn identifier")
		}
		return beginLex
	}
	for _, tok := range looseTokens {
		if l.acceptToken(tok) {
			return beginLex
		}
	}
	return nil
}

func lex(input string) *lexer {
	yyDebug = 3 // 0 is off, 4 is most
	l := &lexer{
		input:  input,
		tokens: make([]Token, 128),
		state:  beginLex,
		stack:  newStack(),
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
