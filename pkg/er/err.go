package er

import (
	"errors"
	"strings"
)

type Err struct {
	op   string
	err  error
	next error
}

func New(reason string, operation string, origin error) *Err {
	return &Err{
		op:   operation,
		err:  errors.New(reason),
		next: origin,
	}
}

func (e *Err) Error() string {
	b := &strings.Builder{}
	writeTo(b, e)

	return b.String()
}

func (e *Err) Unwrap() error { return e.next }

func writeTo(b *strings.Builder, err error) {
	if b.Len() > 0 {
		b.WriteString(" | ")
	}

	e, ok := err.(*Err)
	if !ok {
		b.WriteString(err.Error())
		return
	}

	if e == nil || e.err == nil {
		b.WriteString("nil")
		return
	}

	if e.op != "" {
		b.WriteString(e.op)
		b.WriteString(": ")
	}
	b.WriteString(e.err.Error())

	if e.next != nil {
		writeTo(b, e.next)
	}

}
