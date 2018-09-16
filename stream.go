package caddycfg

import (
	"github.com/mholt/caddy"
)

// Stream stream of tokens with so called confirmation. It works in the following way:
//   1. Next checks if the previous token read was confirmed. Returns true if it is not – Token will return the
// previous token. Looks for the next token and return true if it was found. Otherwise returns false.
//   2. NextArg works like Next but only works for the rest of tokens in the current line
//   3. Token returns token scanned with Next or NextArg
//   4. Confirm confirmes scanned token was consumed
type Stream interface {
	Next() bool
	NextArg() bool
	Token() Token
	Confirm()
}

type streamImpl struct {
	src       *caddy.Controller
	cur       Token
	finished  bool
	confirmed bool
}

// newStream return stream implementation
func newStream(src *caddy.Controller) Stream {
	return &streamImpl{
		src:       src,
		confirmed: true,
	}
}

// Next checks if next token is available
func (t *streamImpl) Next() bool {
	if t.finished {
		return false
	}
	if !t.confirmed {
		return true
	}
	if t.src.Next() {
		t.cur.File = t.src.File()
		t.cur.Value = t.src.Val()
		t.cur.Lin = t.src.Line()
		t.confirmed = false
		return true
	}
	t.finished = true
	return false
}

// NextArg checks if next token at current row is available
func (t *streamImpl) NextArg() bool {
	if t.finished {
		return false
	}
	if !t.confirmed {
		return true
	}
	if t.src.NextArg() {
		t.cur.File = t.src.File()
		t.cur.Value = t.src.Val()
		t.cur.Lin = t.src.Line()
		t.confirmed = false
		return true
	}
	return false
}

// Token ...
func (t *streamImpl) Token() Token {
	return t.cur
}

// Confirm ...
func (t *streamImpl) Confirm() {
	t.confirmed = true
}
