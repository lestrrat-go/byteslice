package byteslice_test

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/lestrrat-go/byteslice"
	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	encoders := []*base64.Encoding{
		base64.RawURLEncoding,
		base64.RawStdEncoding,
		base64.URLEncoding,
		base64.StdEncoding,
	}
	message := []byte("Alice")
	t.Run("Decode", func(t *testing.T) {
		var v byteslice.Type

		var testcases []string
		for _, enc := range encoders {
			testcases = append(testcases, strconv.Quote(enc.EncodeToString(message)))
		}

		for _, tc := range testcases {
			tc := tc
			t.Run(tc, func(t *testing.T) {
				require.NoError(t, json.Unmarshal([]byte(tc), &v), `json.Unmarshal should succeed`)
				require.Equal(t, v.Bytes(), message)
			})
		}
	})
	t.Run("Encode", func(t *testing.T) {
		encmap := make(map[*base64.Encoding]byteslice.Base64Encoder)
		encmap[base64.RawURLEncoding] = byteslice.RawURLEncoder
		encmap[base64.RawStdEncoding] = byteslice.RawStdEncoder
		encmap[base64.URLEncoding] = byteslice.URLEncoder
		encmap[base64.StdEncoding] = byteslice.StdEncoder
		for enc, bsenc := range encmap {
			dst := []byte(strconv.Quote(enc.EncodeToString(message)))
			bsenc := bsenc
			t.Run(string(dst), func(t *testing.T) {
				var v byteslice.Type
				v.SetBytes([]byte("Alice"))
				v.SetEncoder(bsenc)
				buf, err := json.Marshal(v)
				require.NoError(t, err, `json.Marshal should succeed`)
				require.Equal(t, buf, dst, `encoded values should match`)
			})
		}
	})
}

func TestStruct(t *testing.T) {
	var foo struct {
		Bar byteslice.Type `json:"bar"`
	}

	const src = `{"bar": "QWxpY2U"}`

	require.NoError(t, json.Unmarshal([]byte(src), &foo))
	require.Equal(t, foo.Bar.String(), `Alice`)
}
