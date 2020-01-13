package json

import (
	"bytes"
	stdJSON "encoding/json"
	"testing"

	"github.com/eyenih/go-jsonapi/encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testValue struct {
	content string
	t       json.Type
}

func (tv testValue) Type() json.Type {
	return tv.t
}

func (tv *testValue) Read(p []byte) (n int, err error) {
	return 0, nil
}

func TestMarshal(t *testing.T) {
	t.Run("empty values", func(t *testing.T) {
		emptyValue := "{}"
		w := &bytes.Buffer{}
		require.NoError(t, stdJSON.Compact(w, []byte(emptyValue)))
		contentLength := w.Len()
		w.Reset()

		tv := &testValue{t: json.Object}
		n, err := json.Copy(w, tv)
		require.NoError(t, err)

		assert.Equal(t, int64(contentLength), n)
		assert.Equal(t, emptyValue, w.String())
	})
}
