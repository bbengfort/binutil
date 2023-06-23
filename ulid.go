package binutil

import "github.com/oklog/ulid/v2"

type ULID struct {
	ULID ulid.ULID
}

var (
	_ Encoder = &ULID{}
	_ Decoder = &ULID{}
)

func (u ULID) DecodeBinary(in []byte) (_ Encoder, err error) {
	if err = u.ULID.UnmarshalBinary(in); err != nil {
		return nil, err
	}
	return u, nil
}

func (u ULID) DecodeString(in string) (_ Encoder, err error) {
	if u.ULID, err = ulid.ParseStrict(in); err != nil {
		return nil, err
	}
	return u, nil
}

func (u ULID) EncodeBinary() ([]byte, error) {
	return u.ULID.Bytes(), nil
}

func (u ULID) EncodeString() (string, error) {
	return u.ULID.String(), nil
}
