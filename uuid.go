package binutil

import "github.com/google/uuid"

type UUID struct {
	UUID uuid.UUID
}

var (
	_ Encoder = &UUID{}
	_ Decoder = &UUID{}
)

func (u UUID) DecodeBinary(in []byte) (_ Encoder, err error) {
	if err = u.UUID.UnmarshalBinary(in); err != nil {
		return nil, err
	}
	return u, nil
}

func (u UUID) DecodeString(in string) (_ Encoder, err error) {
	if u.UUID, err = uuid.Parse(in); err != nil {
		return nil, err
	}
	return u, nil
}

func (u UUID) EncodeBinary() ([]byte, error) {
	return u.UUID.MarshalBinary()
}

func (u UUID) EncodeString() (string, error) {
	return u.UUID.String(), nil
}
