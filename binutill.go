package binutil

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

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
func New(steps ...any) (_ *Pipeline, err error) {
	decoders := make([]Decoder, 0, len(steps))
	for _, step := range steps {
		var decoder Decoder
		switch t := step.(type) {
		case string:
			if decoder, err = NewDecoder(t); err != nil {
				return nil, err
			}
		case Decoder:
			decoder = t
		default:
			return nil, ErrUnknownStepType
		}
		decoders = append(decoders, decoder)
	}
	return &Pipeline{steps: decoders}, nil
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

var (
	decmu    sync.RWMutex
	decoders map[string]decoder
)

// DecoderConstructors are functions that create a new Decoder ready for use.
type DecoderConstructor func() Decoder

type decoder struct {
	constructor DecoderConstructor
	alias       bool
}

// Register a decoder constructor so that the decoder can be referenced by the name
// suplied and users can instantiate it directly from the type name. Note that names are
// case insensitive so MyDecoder is the same as mydecoder.
func RegisterDecoder(name string, constructor DecoderConstructor, aliases ...string) {
	// All lookups are case insensitive
	name = strings.TrimSpace(strings.ToLower(name))

	decmu.Lock()
	defer decmu.Unlock()
	if decoders == nil {
		decoders = make(map[string]decoder)
	}

	decoders[name] = decoder{constructor, false}
	for _, alias := range aliases {
		decoders[alias] = decoder{constructor, true}
	}
}

// Create a decoder by name rather than by directly instantiating one.
func NewDecoder(name string) (Decoder, error) {
	// All lookups are case insensitive
	name = strings.TrimSpace(strings.ToLower(name))

	decmu.RLock()
	defer decmu.RUnlock()
	if decoder, ok := decoders[name]; ok {
		return decoder.constructor(), nil
	}
	return nil, fmt.Errorf("no registered decoder with the name %q", name)
}

func DecoderNames() []string {
	decmu.RLock()
	defer decmu.RUnlock()

	out := make([]string, 0, len(decoders))
	for name, decoder := range decoders {
		if !decoder.alias {
			out = append(out, name)
		}
	}

	sort.Strings(out)
	return out
}
