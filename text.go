package binutil

import (
	"bytes"
	"io"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
)

// Text encoding schemes to convert to and from.
const (
	UTF8Encoding TextEncoding = iota
	ASCIIEncoding
	Latin1Encoding
)

func init() {
	RegisterDecoder(TextDecoder, func() Decoder { return NewText(UTF8Encoding) })
	RegisterDecoder(UTF8Decoder, func() Decoder { return NewText(UTF8Encoding) })
	RegisterDecoder("utf8", func() Decoder { return NewText(UTF8Encoding) })
	RegisterDecoder("txt", func() Decoder { return NewText(UTF8Encoding) })
	RegisterDecoder(ASCIIDecoder, func() Decoder { return NewText(ASCIIEncoding) })
	RegisterDecoder(Latin1Decoder, func() Decoder { return NewText(Latin1Encoding) })
}

const (
	TextDecoder   = "text"
	UTF8Decoder   = "utf-8"
	ASCIIDecoder  = "ascii"
	Latin1Decoder = "latin1"
)

func NewText(encoding TextEncoding) *Text {
	return &Text{Encoding: encoding}
}

type Text struct {
	Encoding TextEncoding
	data     []byte
}

var (
	_ Encoder = &Text{}
	_ Decoder = &Text{}
)

func (u Text) DecodeBinary(in []byte) (_ Encoder, err error) {
	return &Text{Encoding: u.Encoding, data: in}, nil
}

func (u Text) DecodeString(in string) (_ Encoder, err error) {
	var data []byte
	switch u.Encoding {
	case UTF8Encoding, ASCIIEncoding:
		data = []byte(in)
	default:
		var charset encoding.Encoding
		if charset, err = ianaindex.MIME.Encoding(u.Encoding.String()); err != nil {
			return nil, err
		}

		rx := transform.NewReader(bytes.NewBufferString(in), charset.NewDecoder())
		if data, err = io.ReadAll(rx); err != nil {
			return nil, err
		}
	}

	return &Text{Encoding: UTF8Encoding, data: data}, nil
}

func (u Text) EncodeBinary() ([]byte, error) {
	if u.data != nil {
		return u.data, nil
	}
	return nil, ErrNoData
}

func (u Text) EncodeString() (_ string, err error) {
	switch u.Encoding {
	case UTF8Encoding, ASCIIEncoding:
		return string(u.data), nil
	default:
		var charset encoding.Encoding
		if charset, err = ianaindex.MIME.Encoding(u.Encoding.String()); err != nil {
			return "", err
		}

		var data []byte
		rx := transform.NewReader(bytes.NewBuffer(u.data), charset.NewDecoder())
		if data, err = io.ReadAll(rx); err != nil {
			return "", err
		}
		return string(data), nil
	}
}

type TextEncoding uint8

func (b TextEncoding) String() string {
	switch b {
	case UTF8Encoding:
		return "utf-8"
	case ASCIIEncoding:
		return "ascii"
	case Latin1Encoding:
		return "latin1"
	default:
		return "unknown"
	}
}
