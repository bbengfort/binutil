package binutil

import "encoding/base64"

// Base64 Encoding Schemes for determining the character set and padding used. Standard
// encoding uses the RFC 4648 encoding standard and includes padding characters. URL
// encoding uses the RFC 4648 alternate encoding that is safe for filenames and URLs.
// RawStd and RawURL omits the padding characters.
const (
	B64SchemeStd Base64Scheme = iota
	B64SchemeRawStd
	B64SchemeURL
	B64SchemeRawURL
)

func init() {
	RegisterDecoder(Base64Decoder, func() Decoder { return NewBase64(B64SchemeStd) }, "b64")
	RegisterDecoder(StdBase64Decoder, func() Decoder { return NewBase64(B64SchemeStd) }, "base64std", "b64std")
	RegisterDecoder(RawBase64Decoder, func() Decoder { return NewBase64(B64SchemeRawStd) }, "base64raw", "b64raw")
	RegisterDecoder(URLBase64Decoder, func() Decoder { return NewBase64(B64SchemeURL) }, "base64url", "b64url")
	RegisterDecoder(RawURLBase64Decoder, func() Decoder { return NewBase64(B64SchemeRawURL) }, "base64rawurl", "b64rawurl")
}

const (
	Base64Decoder       = "base64"
	StdBase64Decoder    = "base64-std"
	RawBase64Decoder    = "base64-raw"
	URLBase64Decoder    = "base64-url"
	RawURLBase64Decoder = "base64-rawurl"
)

func NewBase64(scheme Base64Scheme) *Base64 {
	return &Base64{Scheme: scheme}
}

// Base64 implements the encoder and decoder interface for Base64 data and strings.
// Base64 is either an initial decoder or final encoder type and is not used for
// intermediate binary representations.
//
// The primary paramater for base64 is the scheme which determines what character set
// and padding is used in the base64 encoding. Standard encoding is the default.
type Base64 struct {
	Scheme Base64Scheme
	data   []byte
}

var (
	_ Encoder = &Base64{}
	_ Decoder = &Base64{}
)

// DecodeBinary returns a new Base64 object with the wrapped data, ready to be encoded
// as a base64 string.
func (b Base64) DecodeBinary(in []byte) (Encoder, error) {
	return &Base64{Scheme: b.Scheme, data: in}, nil
}

// DecodeString decodes the base64 string with the specified scheme and returns an
// encoder with the data bytes ready to be fetched.
func (b Base64) DecodeString(in string) (_ Encoder, err error) {
	var data []byte
	switch b.Scheme {
	case B64SchemeStd:
		if data, err = base64.StdEncoding.DecodeString(in); err != nil {
			return nil, err
		}
	case B64SchemeRawStd:
		if data, err = base64.RawStdEncoding.DecodeString(in); err != nil {
			return nil, err
		}
	case B64SchemeURL:
		if data, err = base64.URLEncoding.DecodeString(in); err != nil {
			return nil, err
		}
	case B64SchemeRawURL:
		if data, err = base64.RawURLEncoding.DecodeString(in); err != nil {
			return nil, err
		}
	}
	return b.DecodeBinary(data)
}

// EncodeBinary returns the wrapped data if any is available, otherwise returns an error.
func (b Base64) EncodeBinary() ([]byte, error) {
	if b.data != nil {
		return b.data, nil
	}
	return nil, ErrNoData
}

// EncodeString encodes the wrapped data according to the base64 encoding scheme.
func (b Base64) EncodeString() (string, error) {
	if b.data == nil {
		return "", ErrNoData
	}

	switch b.Scheme {
	case B64SchemeStd:
		return base64.StdEncoding.EncodeToString(b.data), nil
	case B64SchemeRawStd:
		return base64.RawStdEncoding.EncodeToString(b.data), nil
	case B64SchemeURL:
		return base64.URLEncoding.EncodeToString(b.data), nil
	case B64SchemeRawURL:
		return base64.RawURLEncoding.EncodeToString(b.data), nil
	default:
		return "", ErrUnknownB64Scheme
	}
}

type Base64Scheme uint8

func (b Base64Scheme) String() string {
	switch b {
	case B64SchemeStd:
		return "StdEncoding"
	case B64SchemeRawStd:
		return "RawStdEncoding"
	case B64SchemeURL:
		return "URLEncoding"
	case B64SchemeRawURL:
		return "RawURLEncoding"
	default:
		return "unknown"
	}
}
