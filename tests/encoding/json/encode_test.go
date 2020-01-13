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
	readErr error
}

func (tv testValue) Type() json.Type {
	return tv.t
}

func (tv *testValue) Read(p []byte) (n int, err error) {
	return 0, tv.readErr
}

func compact(b *bytes.Buffer, json string) int {
	defer b.Reset()

	err := stdJSON.Compact(b, []byte(json))
	if err != nil {
		panic(err)
	}

	return b.Len()
}

func TestEncode(t *testing.T) {
	expectedZeroValues := map[json.Type]string{
		json.Number:  "0",
		json.String:  "\"\"",
		json.Boolean: "false",
		json.Array:   "[]",
		json.Object:  "{}",
	}

	t.Run("zero values", func(t *testing.T) {
		for tp, expectedContent := range expectedZeroValues {
			w := &bytes.Buffer{}
			expectedLength := compact(w, expectedContent)

			tv := &testValue{t: tp}
			e := json.NewEncoder(0)
			n, err := e.Encode(w, tv)
			require.NoError(t, err)

			assert.Equal(t, expectedLength, n)
			assert.Equal(t, expectedContent, w.String())
		}
	})

	t.Run("null values", func(t *testing.T) {
		expectedContent := "null"
		for tp, _ := range expectedZeroValues {
			w := &bytes.Buffer{}

			tv := &testValue{t: tp, readErr: json.ErrValueIsNull}
			e := json.NewEncoder(0)
			n, err := e.Encode(w, tv)
			require.NoError(t, err)

			assert.Equal(t, 4, n)
			assert.Equal(t, expectedContent, w.String())
		}
	})
}
