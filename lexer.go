package abc

import (
	"fmt"
	"strings"
)

type item interface{}

type itemError struct {
	msg string
	ctx string
}

type itemField struct {
	key byte
	val string
}

type itemNote struct {
	acc  string
	note byte
	oct  string
	ryt  string
}

func (n itemNote) String() string {
	return fmt.Sprintf("%s%c%s%s", n.acc, n.note, n.oct, n.ryt)
}

type lexerFunc func(*lexer) lexerFunc

type lexer struct {
	items chan item     // output channel of items
	quit  chan struct{} // quit signal chan
	input string        // text to parse
	start int           // start pos of current item
	pos   int           // current position in the input
}

func (l *lexer) run() {
	for state := lexHeader; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func (l *lexer) stop() {
	close(l.quit)
}

func (l *lexer) emit(i item) bool {
	select {
	case l.items <- i:
		return true
	case <-l.quit:
		return false
	}
}

func (l *lexer) next() (byte, bool) {
	if l.pos >= len(l.input) {
		return 0, false
	}
	l.pos++
	return l.input[l.pos-1], true
}

func (l *lexer) back() *lexer {
	if l.pos > 0 {
		l.pos--
	}
	return l
}

func (l *lexer) current() string {
	res := l.input[l.start:l.pos]
	return res
}

func (l *lexer) advance() {
	l.start = l.pos
}

func (l *lexer) rest() string {
	return l.input[l.pos:]
}

func (l *lexer) toEOL() string {
	rest := l.rest()
	r := strings.Index(rest, "\n")
	if r < 0 {
		l.pos = len(l.input)
	} else {
		l.pos += r + 1
		rest = rest[:r]
	}
	return rest
}

func (l *lexer) errorf(msg string, args ...interface{}) {
	l.emit(itemError{
		msg: fmt.Sprintf(msg, args...),
		ctx: l.current(),
	})
}

func lexHeader(l *lexer) lexerFunc {
	k, ok := l.next()
	if !ok {
		return nil
	}

	colon, ok := l.next()
	if !ok || colon != ':' {
		if ok {
			l.back() // back the colon
		}
		l.back() // back the note
		return lexTune
	}

	ok = l.emit(itemField{
		key: k,
		val: l.toEOL(),
	})

	if !ok {
		return nil
	}
	return lexHeader
}

func isAcc(a byte) bool {
	return a == '^' || a == '_' || a == '='
}

func isNote(n byte) bool {
	return ('a' <= n && n <= 'g') ||
		('A' <= n && n <= 'G') || n == 'z'
}

func isRhythm(n byte) bool {
	return '0' <= n && n <= '9' || n == '/' || n == '>' || n == '<'
}

func isIgnored(c byte) bool {
	return c == ' ' || c == '[' || c == ']'
}

func lexTune(l *lexer) lexerFunc {
	n, ok := l.next()
	for ok && isIgnored(n) {
		l.advance()
		n, ok = l.next()
	}
	if !ok {
		return nil
	}

	var note itemNote

	// accidentals
	if isAcc(n) {
		acc := n
		n, ok = l.next()
		if !ok {
			l.errorf("unfinished note")
			return nil
		}
		if n != acc { // double accidental?
			if n == '^' || n == '_' || n == '=' {
				l.errorf("wrong accidental")
				return nil
			}
			l.back()
		}
		note.acc = l.current()
		l.advance()
		n, ok = l.next()
	}

	// note name
	if !isNote(n) {
		l.errorf("not a note")
	}
	note.note = n
	l.advance()

	// octave modifier
	n, ok = l.next()
	for ok && n == ',' || n == '\'' {
		n, ok = l.next()
	}
	if ok {
		l.back()
	}
	note.oct = l.current()
	l.advance()

	// length
	n, ok = l.next()
	for ok && isRhythm(n) {
		n, ok = l.next()
	}
	if ok {
		l.back()
	}
	note.ryt = l.current()
	l.advance()

	if !l.emit(note) {
		return nil
	}

	return lexTune
}

func lex(input string) (*lexer, chan item) {
	l := &lexer{
		input: input,
		items: make(chan item),
		quit:  make(chan struct{}),
	}
	go l.run()
	return l, l.items
}
