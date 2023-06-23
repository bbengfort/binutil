package binutil

import "errors"

var (
	ErrEmptyPipeline    = errors.New("the pipeline has no transformation steps")
	ErrOverwrite        = errors.New("this operation will overwrite existing data")
	ErrNoData           = errors.New("data cannot be empty or nil")
	ErrUnknownB64Scheme = errors.New("unknown base64 encoding scheme")
	ErrUnknownStepType  = errors.New("initialize a pipeline with a string or Decoder")
)
