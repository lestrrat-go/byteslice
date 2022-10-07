byteslice ![](https://github.com/lestrrat-go/byteslice/workflows/CI/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/lestrrat-go/byteslice.svg)](https://pkg.go.dev/github.com/lestrrat-go/byteslice)
=========

A thin layer on top of `[]byte` to accomodate base64 encoding used in conjunction
with "encoding/json".

When the Go type `[]byte` is used to decode/encode a JSON string, the standard
"encoding/json" library assumes base64.StdEncoding (RFC4648 section 3.2). This
poses a problem when you want to encode the fields using a different encoding,
or if whatever protocol you are using does not very clearly specify which
base64 encoding you are supposed to use and you must field multiple,
slightly differing base64 encoded string.

While the standard "encoding/json" library only works with a single encoding,
if you use `byteslice.Buffer`, it can decode using any of the pre-defined
encodings in "encoding/base64" by default (or specify an alternate `*base64.Encoding`
object), and to specify which encoding to use converting `[]byte` to JSON.

# SYNOPSIS

```go
type Foo struct {
  Field byteslice.Buffer `json:"field"`
}

func init() {
  byteslice.SetGlobalB64Encoder(base64.RawURLEncoding)
  byteslice.SetGlobalB64Decoder(base64.RawURLEncoding)
}

var foo Foo
_ = json.Unmarshal(data, &foo)

```

# FAQ

## Q: What's with `AcceptValue`?

`AcceptValue` is a convenience function for those cases when you do not know
the type of source value before hand, but you still would like to attempt
to initialize a `byteslice.Buffer` object. (This happens more often than you may think!)
