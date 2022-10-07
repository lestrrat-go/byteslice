package byteslice_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/lestrrat-go/byteslice"
	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	encoders := map[string]*base64.Encoding{
		"RawURL": base64.RawURLEncoding,
		"RawStd": base64.RawStdEncoding,
		"URL":    base64.URLEncoding,
		"Std":    base64.StdEncoding,
	}
	message := make([]byte, 64)
	for i := 0; i < 64; i++ {
		message[i] = byte(i)
	}

	t.Run("Decode", func(t *testing.T) {
		var v byteslice.Buffer

		type decodeTC struct {
			Name    string
			Payload string
		}
		var testcases []decodeTC
		for name, enc := range encoders {
			encoded := enc.EncodeToString(message)
			testcases = append(testcases, decodeTC{
				Name:    fmt.Sprintf("%s#%s", name, encoded),
				Payload: strconv.Quote(encoded),
			})
		}

		for _, tc := range testcases {
			tc := tc
			t.Run(tc.Name, func(t *testing.T) {
				require.NoError(t, json.Unmarshal([]byte(tc.Payload), &v), `json.Unmarshal should succeed`)
				require.Equal(t, v.Bytes(), message)
			})
		}
	})
	t.Run("Encode", func(t *testing.T) {
		for _, enc := range encoders {
			dst := []byte(strconv.Quote(enc.EncodeToString(message)))
			enc := enc
			t.Run(string(dst), func(t *testing.T) {
				var v byteslice.Buffer
				v.SetBytes(message)
				v.SetEncoder(enc)
				buf, err := json.Marshal(v)
				require.NoError(t, err, `json.Marshal should succeed`)
				require.Equal(t, buf, dst, `encoded values should match`)
			})
		}
	})
}

func TestStruct(t *testing.T) {
	var foo struct {
		Bar byteslice.Buffer `json:"bar"`
	}

	const src = `{"bar":"QWxpY2U"}`

	require.NoError(t, json.Unmarshal([]byte(src), &foo))
	require.Equal(t, string(foo.Bar.Bytes()), `Alice`)
}
