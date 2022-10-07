package byteslice

import (
	"encoding/base64"
	"strings"
	"sync"
)

// B64Decoder is the interface for objects that can decode
// base64 encoded strings into `[]byte`
//
// Any `*base64.Encoding` object satisfies this interface.
type B64Decoder interface {
	DecodeString(string) ([]byte, error)
}

// B64Encoder is the interface for objects that can encode
// `[]byte` into base64 encoded string
//
// Any `*base64.Encoding` object satisfies this interface.
type B64Encoder interface {
	EncodeToString([]byte) string
}

var globalMu sync.RWMutex
var globalDecoder B64Decoder
var globalEncoder B64Encoder

// SetGlobalB64Decoder sets the `B64Decoder` that should be used globally
func SetGlobalB64Decoder(dec B64Decoder) {
	globalMu.Lock()
	defer globalMu.Unlock()

	globalDecoder = dec
}

// SetGlobalB64Encoder sets the `B64Encoder` that should be used globally
func SetGlobalB64Encoder(enc B64Encoder) {
	globalMu.Lock()
	defer globalMu.Unlock()

	globalEncoder = enc
}

// GlobalB64Decoder returns the `B64Decoder` that is to be used by default
// for all `byteslice.Buffer` types. Each instance can be configured to
// use its own decoder if set individually.
//
// The default decoder uses heuristics to determine which of the `"encoding/base64".Encoding`
// objects should be used.
//
//   - If the incoming payload does NOT contain a '=' at the end of the string, it is considered to be a "raw" base64 encoding (i.e. no padding)
//   - If the incoming payload does NOT contain either '+' or '/', it is considered to be a "url" encoding
//
// By combining these two, we decide which of `base64.URLEncoding`, `base64.StdEncoding`,
// `base64.RawURLEncoding`, or `base64.RawStdEncoding` we should be using to
// decode the JSON string
func GlobalB64Decoder() B64Decoder {
	globalMu.RLock()
	defer globalMu.RUnlock()

	return globalDecoder
}

// GlobalB64Encoder returns the `B64Encoder` that is to be used by default
// for all `byteslice.Buffer` types. Each instance can be configured to
// use its own encoder if set individually.
//
// The default encoder uses the same encoder as the standard library's
// "encoding/json", which is the `base64.StdEncoding`
func GlobalB64Encoder() B64Encoder {
	globalMu.RLock()
	defer globalMu.RUnlock()

	return globalEncoder
}

// B64DecoderFunc is an instance of B64Decoder that is based on
// a function.
type B64DecoderFunc func(string) ([]byte, error)

// DecodeString implements the B64Decoder interface
func (f B64DecoderFunc) DecodeString(s string) ([]byte, error) {
	return f(s)
}

// B64EncoderFunc is an instance of B64Encoder that is based on
// a function.
type B64EncoderFunc func([]byte) string

// EncodeString implements the B64Encoder interface
func (f B64EncoderFunc) EncodeToString(data []byte) string {
	return f(data)
}

func defaultDecodeString(src string) ([]byte, error) {
	var enc *base64.Encoding

	var isRaw = !strings.HasSuffix(src, "=")
	var isURL = !strings.ContainsAny(src, "+/")
	switch {
	case isRaw && isURL:
		enc = base64.RawURLEncoding
	case isURL:
		enc = base64.URLEncoding
	case isRaw:
		enc = base64.RawStdEncoding
	default:
		enc = base64.StdEncoding
	}

	return enc.DecodeString(src)
}

func init() {
	SetGlobalB64Decoder(B64DecoderFunc(defaultDecodeString))
	SetGlobalB64Encoder(base64.StdEncoding)
}
