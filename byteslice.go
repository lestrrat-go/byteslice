// Package byteslice provides a thin abstraction over `[]byte`
// types. This package exist to support encoding/decoding
// `[]byte` slices in JSON serialization/deserialization.
//
// Yes, `encoding/json` supports base64 encoding for `[]byte`
// fields, but there's no way to customize the way these
// fields get serialized/deserialized -- e.g. with padding
// or no padding, which characters to use for padding, etc.
//
// By using byteslice.Type as the field instead of `[]byte`
// you can change the the way this base64 handling is performed.
package byteslice

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Type struct {
	mu      sync.RWMutex
	data    []byte
	decoder Base64Decoder
	encoder Base64Encoder
}

// B64Decoder returns the Base64Decoder associated with this object.
// If uninitialized, will use the global decoder via byteslice.GlobalB64Decoder()
func (t *Type) B64Decoder() Base64Decoder {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.b64DecoderNoLock()
}

func (t *Type) b64DecoderNoLock() Base64Decoder {
	if t.decoder != nil {
		return t.decoder
	}
	return GlobalB64Decoder()
}

func (t *Type) SetDecoder(dec Base64Decoder) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.decoder = dec
}

// B64Encoder returns the Base64Encoder associated with this object.
// If uninitialized, will use the global decoder via byteslice.GlobalB64Encoder()
func (t *Type) B64Encoder() Base64Encoder {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.b64EncoderNoLock()
}

func (t *Type) b64EncoderNoLock() Base64Encoder {
	if t.encoder != nil {
		return t.encoder
	}
	return GlobalB64Encoder()
}

func (t *Type) SetEncoder(enc Base64Encoder) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.encoder = enc
}

func (t *Type) UnmarshalJSON(data []byte) error {
	if t == nil {
		return fmt.Errorf(`nil byteslice.Type`)
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf(`failed to unmarshal data to byteslice.Type: %w`, err)
	}

	buf, err := t.b64DecoderNoLock().DecodeString(raw)
	if err != nil {
		return fmt.Errorf(`failed to base64 decode unmarshaled data  for byteslice.Type: %w`, err)
	}
	t.data = buf
	return nil
}

func (t Type) MarshalJSON() ([]byte, error) {
	s, err := t.b64EncoderNoLock().EncodeToString(t.data)
	if err != nil {
		return nil, fmt.Errorf(`failed to encode data into base64 string`)
	}
	return json.Marshal(s)
}

func (t *Type) Bytes() []byte {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.data
}

func (t *Type) String() string {
	return string(t.Bytes())
}

func (t *Type) SetBytes(data []byte) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.data = data
}
