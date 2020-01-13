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

func compact(b *bytes.Buffer, json string) int {
	defer b.Reset()

	err := stdJSON.Compact(b, []byte(json))
	if err != nil {
		panic(err)
	}

	return b.Len()
}

func TestEncode(t *testing.T) {
	t.Run("zero values", func(t *testing.T) {
		zeroValues := map[json.Type]string{
			json.Number:  "0",
			json.String:  "\"\"",
			json.Boolean: "false",
			json.Array:   "[]",
			json.Object:  "{}",
		}

		for tp, expectedContent := range zeroValues {
			w := &bytes.Buffer{}
			expectedLength := compact(w, expectedContent)

			tv := &testValue{t: tp}
			e := json.NewEncoder(0)
			n, err := e.Encode(w, tv)
			require.NoError(t, err)

			assert.Equal(t, int64(expectedLength), n)
			assert.Equal(t, expectedContent, w.String())
		}

	})
}
