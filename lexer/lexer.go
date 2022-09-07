package lexer

import (
	"bufio"
	"io"
	"strings"

	"monkey-book/token"
)

type Lexer struct {
	reader       *bufio.Reader
	ch           rune // current char under examination
	lineNumber   int
	columnNumber int
}

func New(reader io.Reader) *Lexer {
	l := &Lexer{reader: bufio.NewReader(reader), lineNumber: 1}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	r, _, err := l.reader.ReadRune()
	if err == io.EOF {
		l.ch = 0
		return
	}
	if l.ch == '\n' {
		l.lineNumber += 1
		l.columnNumber = 1
	} else {
		l.columnNumber += 1
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
				Type:         token.EQ,
				Literal:      buf.String(),
				LineNumber:   l.lineNumber,
				ColumnNumber: l.columnNumber,
			}
		} else {
			tok = l.newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = l.newToken(token.PLUS, l.ch)
	case '-':
		tok = l.newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			var buf strings.Builder
			buf.WriteRune(l.ch)
			l.readChar()
			buf.WriteRune(l.ch)
			tok = token.Token{
				Type:         token.NOT_EQ,
				Literal:      buf.String(),
				LineNumber:   l.lineNumber,
				ColumnNumber: l.columnNumber,
			}
		} else {
			tok = l.newToken(token.BANG, l.ch)
		}
	case '/':
		tok = l.newToken(token.SLASH, l.ch)
	case '*':
		tok = l.newToken(token.ASTERISK, l.ch)
	case '<':
		tok = l.newToken(token.LT, l.ch)
	case '>':
		tok = l.newToken(token.GT, l.ch)
	case ';':
		tok = l.newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = l.newToken(token.LPAREN, l.ch)
	case ')':
		tok = l.newToken(token.RPAREN, l.ch)
	case ',':
		tok = l.newToken(token.COMMA, l.ch)
	case '{':
		tok = l.newToken(token.LBRACE, l.ch)
	case '}':
		tok = l.newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.LineNumber = l.lineNumber
		tok.ColumnNumber = l.columnNumber
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.LineNumber = l.lineNumber
			tok.ColumnNumber = l.columnNumber
			return tok
		} else if isDigit(l.ch) {
			tok = l.readNumber()
			tok.LineNumber = l.lineNumber
			tok.ColumnNumber = l.columnNumber
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
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

func (l *Lexer) readNumber() token.Token {
	var seenDot = false
	var buf strings.Builder
	for {
		if l.ch == '.' {
			if seenDot {
				return l.newToken(token.ILLEGAL, l.ch)
			}
			seenDot = true
		} else if !isDigit(l.ch) {
			break
		}
		buf.WriteRune(l.ch)
		l.readChar()
	}
	var tokenType token.TokenType
	if seenDot {
		tokenType = token.FLOAT
	} else {
		tokenType = token.INT
	}
	return token.Token{
		Type:         tokenType,
		Literal:      buf.String(),
		LineNumber:   l.lineNumber,
		ColumnNumber: l.columnNumber,
	}
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), LineNumber: l.lineNumber, ColumnNumber: l.columnNumber}
}
