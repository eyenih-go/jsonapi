package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/eyenih/go-jsonapi"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testWriterTo struct {
	content string
}

func (wt *testWriterTo) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprint(w, wt.content)

	return int64(n), err
}

func TestDocument(t *testing.T) {
	w := &bytes.Buffer{}
	json.Compact(w, []byte(`{
		"data": null
	}`))
	contentLength := w.Len()
	expected, err := ioutil.ReadAll(w)
	require.NoError(t, err)

	w.Reset()

	sd, err := jsonapi.NewStandardDocument(&testWriterTo{"null"}, nil, nil, nil, nil, nil)
	require.NoError(t, err)

	n, err := sd.WriteTo(w)

	require.NoError(t, err)
	assert.Equal(t, n, int64(contentLength))
	assert.Equal(t, string(expected), w.String())
}
