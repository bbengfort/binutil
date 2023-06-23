package binutil

import "encoding/hex"

type Hex struct {
	data []byte
}

var (
	_ Encoder = &Hex{}
	_ Decoder = &Hex{}
)

func (h Hex) DecodeBinary(in []byte) (Encoder, error) {
	return &Hex{data: in}, nil
}

func (h Hex) DecodeString(in string) (_ Encoder, err error) {
	var data []byte
	if data, err = hex.DecodeString(in); err != nil {
		return nil, err
	}
	return h.DecodeBinary(data)
}

func (h Hex) EncodeBinary() ([]byte, error) {
	if h.data != nil {
		return h.data, nil
	}
	return nil, ErrNoData
}

func (h Hex) EncodeString() (string, error) {
	if h.data == nil {
		return "", ErrNoData
	}
	return hex.EncodeToString(h.data), nil
}
