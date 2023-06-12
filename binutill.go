package binutil

// An encoder is an object that can marshal either a binary or a string representation
// of it's underlying data. For example a UUID is 16 bytes binary or it can be a GUID
// string value. Some representations are strings only; for example JSON is a string
// representation of data and the binary encoding is simply UTF-8 bytes. In other cases
// data is binary only; for example protocol buffers are binary data and the string
// representation may be base64 encoded bytes.
type Encoder interface {
	EncodeBinary() (data []byte, err error)
	EncodeString() (data string, err error)
}

// Decoder is an object that can unmarshal itself from either a binary or string
// representation and in both cases ensure that complete data is returned.
type Decoder interface {
	DecodeBinary(data []byte) (err error)
	DecodeString(data string) (err error)
}

// Transformer converts data from the encoder type to the decoder type.
type Transformer interface {
	Decoder
	Encoder
}
