// Package byteslice provides a thin abstraction over `[]byte`
// types. This package exist to support encoding/decoding
// `[]byte` slices in JSON serialization/deserialization.
//
// Yes, `encoding/json` supports base64 encoding for `[]byte`
// fields, but there's no way to customize the way these
// fields get serialized/deserialized -- e.g. with padding
// or no padding, which characters to use for padding, etc.
//
// By using byteslice.Buffer as the field instead of `[]byte`
// you can change the the way this base64 handling is performed.
package byteslice

import (
	"encoding/json"
	"fmt"
)

// Buffer represents a byte slice. Its only purpose is to act
// as a thing layer on top of a `[]byte` and provide flexibly
// JSON serialization/deserialization capabilities
//
// It is safe to use the zero value of the `Buffer` object,
// but the object is not explicitly synchronized. The user
// must make sure to apply any synchronization if need be.
//
// You should not copy a `Buffer` object by reference
type Buffer struct {
	data    []byte
	decoder B64Decoder
	encoder B64Encoder
}

// New creates a new buffer. Using the data provided as the initial buffer.
// This is different from using `SetBytes`, which copies the values onto
// the internal buffer.
//
// You may pass `nil` to the argument to create an uninitialized `Buffer` object.
//
// If you do not need explicit initialization, it is safe to use the
// zero value of the `Buffer` object.
func New(data []byte) *Buffer {
	return &Buffer{data: data}
}

// B64Decoder returns the B64Decoder associated with this object.
// If uninitialized, it will use the global decoder via byteslice.GlobalB64Decoder()
func (b *Buffer) B64Decoder() B64Decoder {
	if b.decoder != nil {
		return b.decoder
	}
	return GlobalB64Decoder()
}

// SetB64Decoder assigns a B64Decoder for this object.
func (b *Buffer) SetB64Decoder(dec B64Decoder) *Buffer {
	b.decoder = dec
	return b
}

// B64Encoder returns the B64Encoder associated with this object.
// If uninitialized, will use the global decoder via byteslice.GlobalB64Encoder()
func (b *Buffer) B64Encoder() B64Encoder {
	if b.encoder != nil {
		return b.encoder
	}
	return GlobalB64Encoder()
}

// SetB64Encoder assigns a B64Encoder for this object.
func (b *Buffer) SetEncoder(enc B64Encoder) *Buffer {
	b.encoder = enc
	return b
}

// UnmarshalJSON implements `"encoding/json".Unmarshaler`, and provides
// a method to deserialize a `[]byte` string from a base64 encoded
// JSON string.
//
// The JSON string will be parsed using the B64Decoder object associated
// with this object (or the global one, if not specified).
func (b *Buffer) UnmarshalJSON(data []byte) error {
	if b == nil {
		return fmt.Errorf(`nil byteslice.Buffer`)
	}

	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf(`failed to unmarshal data to byteslice.Buffer: %w`, err)
	}

	if err := b.decodeAndSetString(raw); err != nil {
		return fmt.Errorf(`failed to accept unmarshaled data: %w`, err)
	}
	return nil
}

func (b *Buffer) decodeAndSetString(in string) error {
	buf, err := b.B64Decoder().DecodeString(in)
	if err != nil {
		return fmt.Errorf(`failed to decode string for byteslice.Buffer: %w`, err)
	}
	b.data = buf
	return nil
}

// MarshalJSON implements `"encoding/json".Marshaler, and provides
// a method to serialize a `[]byte` string to a base64 encoded
// JSON string.
//
// The JSON string will be parsed using the B64Encoder object associated
// with this object (or the global one, if not specified).
func (b Buffer) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.B64Encoder().EncodeToString(b.data))
}

// Bytes returns the raw bytes stored in the `Buffer` object.
//
// Users need to take care of synchronization or acting upon on the
// returned buffer, as it will affect the actual stored `[]byte` field
// in the `Buffer` object.
func (b *Buffer) Bytes() []byte {
	if b == nil {
		return nil
	}
	return b.data
}

// AcceptValue is used in by some consumers to assign the value
// whose type is not known before hand.
//
// Values can be either one of the following types: `*byteslice.Buffer`,
// `[]byte`, or `string`.
//
// If the value is a `*byteslice.Buffer`, a copy of the underlying
// is created, and assigned to receiver.
//
// If the value is a `[]byte`, it is the same as calling `SetBytes()`
//
// IF the value is a `string`, the string is assumed to be a base64-encoded
// string. Unlike in the case of `UnmarshalJSON`, the string does not need
// to be quoted.
func (b *Buffer) AcceptValue(in interface{}) error {
	switch in := in.(type) {
	case *Buffer:
		b.SetBytes(in.Bytes())
		return nil
	case []byte:
		b.SetBytes(in)
		return nil
	case string:
		if err := b.decodeAndSetString(in); err != nil {
			return fmt.Errorf(`failed to accept value for byteslice.Buffer: %w`, err)
		}
		return nil
	default:
		return fmt.Errorf(`failed to accept value for byteslice.Buffer: can't handle type %T`, in)
	}
}

// SetBytes copies the `data` byte slice to the internal buffer
func (b *Buffer) SetBytes(data []byte) {
	l := len(data)
	if cap(b.data) < l {
		b.data = make([]byte, l)
	} else {
		b.data = b.data[:l]
	}
	copy(b.data, data)
}

// Len returns the length of the internal `[]byte` buffer
func (b *Buffer) Len() int {
	if b == nil {
		return 0
	}
	return len(b.data)
}
