package binutil

import "fmt"

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
	DecodeBinary(data []byte) (Encoder, error)
	DecodeString(data string) (Encoder, error)
}

// Pipelines manage transformers converting data from the input type to the output type.
type Pipeline struct {
	steps []Decoder
}

// New returns a pipeline that can convert binary and string data.
func New(steps ...Decoder) *Pipeline {
	return &Pipeline{steps: steps}
}

// Bin2Bin transforms binary input data into binary output data by decoding the binary
// data at each step of the pipeline and encoding it to bytes before passing it to the
// next step in the pipeline.
func (p *Pipeline) Bin2Bin(in []byte) (_ []byte, err error) {
	if len(p.steps) == 0 {
		return nil, ErrEmptyPipeline
	}

	for i, step := range p.steps {
		var encoder Encoder
		if encoder, err = step.DecodeBinary(in); err != nil {
			return nil, fmt.Errorf("could not decode binary in step %d: %w", i, err)
		}

		if in, err = encoder.EncodeBinary(); err != nil {
			return nil, fmt.Errorf("could not encode binary in step %d: %w", i, err)
		}
	}

	return in, nil
}

// Bin2Str transforms binary input data into a string representation by decoding the
// binary input data at each step of the pipeline and encoding it to bytes before
// passing it to the next step in the pipeline. The final step is encoded as a str.
func (p *Pipeline) Bin2Str(in []byte) (out string, err error) {
	if len(p.steps) == 0 {
		return "", ErrEmptyPipeline
	}

	lastStep := len(p.steps) - 1
	for i, step := range p.steps {
		var encoder Encoder
		if encoder, err = step.DecodeBinary(in); err != nil {
			return "", fmt.Errorf("could not decode binary in step %d: %w", i, err)
		}

		if i == lastStep {
			if out, err = encoder.EncodeString(); err != nil {
				return "", fmt.Errorf("could not encode string in step %d: %w", i, err)
			}
		} else {
			if in, err = encoder.EncodeBinary(); err != nil {
				return "", fmt.Errorf("could not encode binary in step %d: %w", i, err)
			}
		}
	}

	return out, nil
}

// Str2Bin transforms binary input data into binary output data by decoding the string
// in the first step of the pipeline then encoding it to bytes before passing it to each
// additional step to decode as bytes.
func (p *Pipeline) Str2Bin(in string) (out []byte, err error) {
	if len(p.steps) == 0 {
		return nil, ErrEmptyPipeline
	}

	var encoder Encoder
	if encoder, err = p.steps[0].DecodeString(in); err != nil {
		return nil, fmt.Errorf("could not decode string in step %d: %w", 0, err)
	}

	if out, err = encoder.EncodeBinary(); err != nil {
		return nil, fmt.Errorf("could not encode binary in step %d: %w", 0, err)
	}

	if len(p.steps) > 1 {
		for i, step := range p.steps[1:] {
			if encoder, err = step.DecodeBinary(out); err != nil {
				return nil, fmt.Errorf("could not decode binary in step %d: %w", i, err)
			}

			if out, err = encoder.EncodeBinary(); err != nil {
				return nil, fmt.Errorf("could not encode binary in step %d: %w", i, err)
			}
		}
	}

	return out, nil
}

// Str2Str transforms string input data into a different string representation by
// decoding the string input data at the first step of the pipeline then encoding it to
// bytes and decoding as binary for each additional step of the pipeline The final step
// is encoded as a string.
func (p *Pipeline) Str2Str(in string) (out string, err error) {
	if len(p.steps) == 0 {
		return "", ErrEmptyPipeline
	}

	var encoder Encoder
	if encoder, err = p.steps[0].DecodeString(in); err != nil {
		return "", fmt.Errorf("could not decode string in step %d: %w", 0, err)
	}

	if len(p.steps) > 1 {
		var data []byte
		if data, err = encoder.EncodeBinary(); err != nil {
			return "", fmt.Errorf("could not encode binary in step %d: %w", 0, err)
		}

		lastStep := len(p.steps) - 2
		for i, step := range p.steps[1:] {
			if encoder, err = step.DecodeBinary(data); err != nil {
				return "", fmt.Errorf("could not decode binary in step %d: %w", i, err)
			}

			if i == lastStep {
				if out, err = encoder.EncodeString(); err != nil {
					return "", fmt.Errorf("could not encode string in step %d: %w", i, err)
				}
			} else {
				if data, err = encoder.EncodeBinary(); err != nil {
					return "", fmt.Errorf("could not encode binary in step %d: %w", i, err)
				}
			}
		}
	} else {
		if out, err = encoder.EncodeString(); err != nil {
			return "", fmt.Errorf("could not encode string in step %d: %w", 0, err)
		}
	}

	return out, nil
}
