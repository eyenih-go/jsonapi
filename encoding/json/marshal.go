package json

import (
	"errors"
	"io"
)

var (
	ErrEmptyBuffer = errors.New("Buffer size can't be 0.")
	ErrValueIsNull = errors.New("Value is null.")
)

const (
	Number Type = iota + 1
	String
	Boolean
	Array
	Object
)

const Null = "null"

var emptyValues = map[Type]string{
	Number:  "0",
	String:  "\"\"",
	Boolean: "false",
	Array:   "[]",
	Object:  "{}",
}

type Type int8

type Value interface {
	io.Reader
	Type() Type
}

type ArrayValue interface {
	Values() []Value
}

type ObjectValue interface {
	Values() map[string]Value
}

type Encoder struct {
	buf []byte
}

func NewEncoder(bufferSize int) *Encoder {
	if bufferSize == 0 {
		bufferSize = 32 * 1024
	}

	return &Encoder{
		buf: make([]byte, bufferSize),
	}
}

func write(w io.Writer, buf []byte, v string) (written int64, err error) {
	copy(buf, v)

	n32, ew := w.Write(buf[:len(v)])
	written += int64(n32)
	return written, ew
}

func writeEmpty(w io.Writer, t Type, buf []byte) (written int64, err error) {
	emptyValue := emptyValues[t]

	if len(buf) < len(emptyValue) {
		return 0, io.ErrShortBuffer
	}

	return write(w, buf, emptyValue)
}

func writeNull(w io.Writer, buf []byte) (written int64, err error) {
	if len(buf) < len(Null) {
		return 0, io.ErrShortBuffer
	}

	return write(w, buf, Null)
}

// If it can't read anything from the given Value (length of all the written bytes are 0 and the currently read bytes
// are 0) then value is a zero-value, a null, or it should omit the value.
// If there was no error during the reading, then it's a zero-value and the proper literals will be written to the Writer.
// If ErrValueIsNull has been returned then "null" literal will be written to the given Writer.
// If io.EOF is the error then it won't write anything to the Writer.
func (e *Encoder) Encode(w io.Writer, v Value) (written int64, err error) {
	if e.buf == nil || len(e.buf) == 0 {
		return 0, ErrEmptyBuffer
	}

	for read, rErr := v.Read(e.buf); nil == rErr; read, rErr = v.Read(e.buf) {
		if written == 0 && read == 0 {
			// there was nothing to read.
			if rErr == nil {
				return writeEmpty(w, v.Type(), e.buf)
			} else if rErr == ErrValueIsNull {
				return writeNull(w, e.buf)
			} else if rErr == io.EOF {
				return 0, rErr
			}

			return 0, rErr
		}
	}

	/*
		for {
			read, rErr := v.Read(e.buf)

			if read > 0 {
				nw, ew := w.Write(e.buf[:read])
				if nw > 0 {
					written += int64(nw)
				}
				if ew != nil {
					err = ew
					break
				}
				if read != nw {
					err = io.ErrShortWrite
					break
				}
			}

			if rErr != nil {
				if rErr != io.EOF {
					err = rErr
				}
				break
			}
		}
	*/
	return written, err
}
