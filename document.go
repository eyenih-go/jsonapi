package jsonapi

import (
	"errors"
	"fmt"
	"io"
)

type StandardDocument struct {
	isData bool

	dataOrErr io.WriterTo
	meta      io.WriterTo
	jsonapi   io.WriterTo
	links     io.WriterTo
	included  io.WriterTo
}

func NewStandardDocument(data, err, meta, jsonapi, links, included io.WriterTo) (*StandardDocument, error) {
	if data == nil && err == nil && meta == nil {
		return nil, errors.New("TODO it should contain at least one")
	}

	if data != nil && err != nil {
		return nil, errors.New("TODO it should contain one of the above")
	}

	sd := &StandardDocument{}
	if data != nil {
		sd.isData = true
		sd.dataOrErr = data
		sd.included = included
	} else {
		sd.dataOrErr = err
	}
	sd.meta = meta
	sd.jsonapi = jsonapi
	sd.links = links

	return sd, nil

}

type omitemptyWriter struct {
	w io.Writer

	key              string
	bytesAboutToCome bool
}

func (ow *omitemptyWriter) Write(p []byte) (sum int, err error) {
	if !ow.bytesAboutToCome {
		n, err := ow.w.Write([]byte(fmt.Sprintf("\"%s\":", ow.key)))

		sum += n
		if err != nil {
			return sum, err
		}

		ow.bytesAboutToCome = true
	}

	n, err := ow.w.Write(p)

	sum += n

	return
}

func (w *omitemptyWriter) refresh() {
	w.bytesAboutToCome = false
}

func (sd StandardDocument) WriteTo(w io.Writer) (sum int64, err error) {
	ow := &omitemptyWriter{w: w}

	n32, err := ow.w.Write([]byte("{"))
	sum += int64(n32)
	if sd.isData {
		ow.key = "data"
	} else {
		ow.key = "errors"
	}

	n64, err := sd.dataOrErr.WriteTo(ow)
	sum += n64

	if err != nil {
		return
	}

	ow.refresh()

	n32, err = ow.w.Write([]byte("}"))

	sum += int64(n32)

	return
}
