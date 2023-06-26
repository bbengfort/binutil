package binutil

import "fmt"

func NewMulti(steps ...string) (_ *MultiPipeline, err error) {
	pipes := make(map[string]*Pipeline, len(steps))
	for _, step := range steps {
		var pipe *Pipeline
		if pipe, err = New(step); err != nil {
			return nil, err
		}
		pipes[step] = pipe
	}
	return &MultiPipeline{steps: pipes}, nil
}

// MultiPipeline is able to convert input data type to multiple output types.
type MultiPipeline struct {
	steps map[string]*Pipeline
}

// Bin2Bin transforms binary input data into binary output data for the named step.
func (p *MultiPipeline) Bin2Bin(name string, in []byte) (_ []byte, err error) {
	pipe, ok := p.steps[name]
	if !ok {
		return nil, fmt.Errorf("no pipeline named %q", name)
	}
	return pipe.Bin2Bin(in)
}

// Bin2Str transforms binary input data into a string representation for the named step.
func (p *MultiPipeline) Bin2Str(name string, in []byte) (out string, err error) {
	pipe, ok := p.steps[name]
	if !ok {
		return "", fmt.Errorf("no pipeline named %q", name)
	}
	return pipe.Bin2Str(in)
}

// Str2Bin transforms binary input data into binary output data for the named step.
func (p *MultiPipeline) Str2Bin(name, in string) (out []byte, err error) {
	pipe, ok := p.steps[name]
	if !ok {
		return nil, fmt.Errorf("no pipeline named %q", name)
	}
	return pipe.Str2Bin(in)
}

// Str2Str transforms string input data into a different string representation for the named step.
func (p *MultiPipeline) Str2Str(name, in string) (out string, err error) {
	pipe, ok := p.steps[name]
	if !ok {
		return "", fmt.Errorf("no pipeline named %q", name)
	}
	return pipe.Str2Str(in)
}

func (p *MultiPipeline) MustBin2Str(name string, in []byte) string {
	out, err := p.Bin2Str(name, in)
	if err != nil {
		panic(err)
	}
	return out
}
