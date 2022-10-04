package byteslice

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"sync"
)

// Base64Decoder is the interface for objects that can decode
// base64 encoded strings into `[]byte`
type Base64Decoder interface {
	DecodeString(string) ([]byte, error)
}

// Base64Encoder is the interface for objects that can encode
// `[]byte` into base64 encoded string
type Base64Encoder interface {
	EncodeToString([]byte) (string, error)
}

var globalMu sync.RWMutex
var globalDecoder Base64Decoder
var globalEncoder Base64Encoder

func SetGlobalB64Decoder(dec Base64Decoder) {
	globalMu.Lock()
	defer globalMu.Unlock()

	globalDecoder = dec
}

func SetGlobalB64Encoder(enc Base64Encoder) {
	globalMu.Lock()
	defer globalMu.Unlock()

	globalEncoder = enc
}

func GlobalB64Decoder() Base64Decoder {
	globalMu.RLock()
	defer globalMu.RUnlock()

	return globalDecoder
}

func GlobalB64Encoder() Base64Encoder {
	globalMu.RLock()
	defer globalMu.RUnlock()

	return globalEncoder
}

type B64DecoderFunc func(string) ([]byte, error)

func (f B64DecoderFunc) DecodeString(s string) ([]byte, error) {
	return f(s)
}

type B64EncoderFunc func([]byte) (string, error)

func (f B64EncoderFunc) EncodeToString(data []byte) (string, error) {
	return f(data)
}

var DefaultB64Decoder Base64Decoder = B64DecoderFunc(defaultDecodeString)
var RawURLDecoder Base64Decoder = B64DecoderFunc(func(src string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(src)
})
var URLDecoder Base64Decoder = B64DecoderFunc(func(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
})
var RawStdDecoder Base64Decoder = B64DecoderFunc(func(src string) ([]byte, error) {
	return base64.RawStdEncoding.DecodeString(src)
})
var StdDecoder Base64Decoder = B64DecoderFunc(func(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
})

var DefaultB64Encoder Base64Encoder = RawURLEncoder
var RawURLEncoder Base64Encoder = B64EncoderFunc(func(src []byte) (string, error) {
	return base64.RawURLEncoding.EncodeToString(src), nil
})
var URLEncoder Base64Encoder = B64EncoderFunc(func(src []byte) (string, error) {
	return base64.URLEncoding.EncodeToString(src), nil
})
var RawStdEncoder Base64Encoder = B64EncoderFunc(func(src []byte) (string, error) {
	return base64.RawStdEncoding.EncodeToString(src), nil
})
var StdEncoder Base64Encoder = B64EncoderFunc(func(src []byte) (string, error) {
	return base64.StdEncoding.EncodeToString(src), nil
})

func defaultDecodeString(src string) ([]byte, error) {
	return defaultDecode([]byte(src))
}

func defaultDecode(src []byte) ([]byte, error) {
	var enc *base64.Encoding

	var isRaw = !bytes.HasSuffix(src, []byte{'='})
	var isURL = !bytes.ContainsAny(src, "+/")
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

	dst := make([]byte, enc.DecodedLen(len(src)))
	n, err := enc.Decode(dst, src)
	if err != nil {
		return nil, fmt.Errorf(`failed to decode source: %w`, err)
	}
	return dst[:n], nil
}

func init() {
	SetGlobalB64Decoder(DefaultB64Decoder)
	SetGlobalB64Encoder(DefaultB64Encoder)
}
