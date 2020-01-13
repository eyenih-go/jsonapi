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
	Integer Type = iota + 1
	Float
	String
	Boolean
	Array
	Object
)

const Null = "null"

var emptyValues = map[Type]string{
	Integer: "0",
	Float:   "0.0",
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

func Copy(w io.Writer, v Value) (sum int64, err error) {
	return CopyBuffer(w, v, make([]byte, 32*1024))
}

func CopyBuffer(w io.Writer, v Value, buf []byte) (n int64, err error) {
	if buf != nil && len(buf) == 0 {
		return 0, ErrEmptyBuffer
	}

	for {
		nr, er := v.Read(buf)
		if n == 0 && nr == 0 {
			// there was nothing to read.
			value := ""
			length := 0

			if er == nil {
				// without error we give back the proper empty value.
				emptyValue := emptyValues[v.Type()]
				length = len(emptyValue)
				if len(buf) < length {
					return 0, io.ErrShortBuffer
				}
				value = emptyValue
			} else {
				if er == ErrValueIsNull {
					value = Null
					length = len(Null)

					if len(buf) < length {
						return 0, io.ErrShortBuffer
					}
				}

				if er == io.EOF {
					return 0, nil
				}

				return 0, er

			}
			copy(buf, value)
			n32, ew := w.Write(buf[0:length])
			n += int64(n32)

			return n, ew
		}

		if nr > 0 {
			nw, ew := w.Write(buf[0:nr])
			if nw > 0 {
				n += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}

		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return n, err
}
