package lexer

import (
	"bufio"
	"io"
	"strings"

	"monkey-book/token"
)

type Lexer struct {
	reader *bufio.Reader
	ch     rune // current char under examination
}

func New(reader io.Reader) *Lexer {
	l := &Lexer{reader: bufio.NewReader(reader)}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	r, _, err := l.reader.ReadRune()
	if err == io.EOF {
		l.ch = 0
		return
	}
	l.ch = r
}

func (l *Lexer) peekChar() rune {
	r, _, err := l.reader.ReadRune()
	if err == io.EOF {
		return 0
	}
	err2 := l.reader.UnreadRune()
	if err2 != nil {
		return 0
	}
	return r
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			var buf strings.Builder
			buf.WriteRune(l.ch)
			l.readChar()
			buf.WriteRune(l.ch)
			tok = token.Token{
				Type:    token.EQ,
				Literal: buf.String(),
			}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			var buf strings.Builder
			buf.WriteRune(l.ch)
			l.readChar()
			buf.WriteRune(l.ch)
			tok = token.Token{
				Type:    token.NOT_EQ,
				Literal: buf.String(),
			}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	var buf strings.Builder
	for isLetter(l.ch) {
		buf.WriteRune(l.ch)
		l.readChar()
	}
	return buf.String()
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	var buf strings.Builder
	for isDigit(l.ch) {
		buf.WriteRune(l.ch)
		l.readChar()
	}
	return buf.String()
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
