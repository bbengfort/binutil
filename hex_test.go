package binutil_test

import (
	"testing"

	"github.com/bbengfort/binutil"
	"github.com/stretchr/testify/require"
)

func TestHex(t *testing.T) {
	// Create encoder and fixtures for tests
	h := binutil.Hex{}
	fixtures := fixtures()

	for _, fixture := range fixtures {
		eb, err := h.DecodeBinary(fixture.data)
		require.NoError(t, err, "could not decode binary for fixture %q", fixture.name)

		data, err := eb.EncodeBinary()
		require.NoError(t, err, "could not encode binary for fixture %q", fixture.name)
		require.Equal(t, fixture.data, data, "expected unchanged binary data for fixture %q", fixture.name)

		s, err := eb.EncodeString()
		require.NoError(t, err, "could not encode string for fixture %q", fixture.name)

		if fixture.name == "empty, non-nil data" {
			require.Empty(t, s, "expected no string to be returned for fixture %q", fixture.name)
		} else {
			require.NotEmpty(t, s, "expected a string to be returned for fixture %q", fixture.name)
		}

		es, err := h.DecodeString(s)
		require.NoError(t, err, "could not decode string for fixture %q", fixture.name)

		data, err = es.EncodeBinary()
		require.NoError(t, err, "could not encode binary from string for fixture %q", fixture.name)
		require.Equal(t, fixture.data, data, "expected unchanged binary data from decoded string for fixture %q", fixture.name)
	}
}

func TestRegisteredHex(t *testing.T) {
	dec, err := binutil.NewDecoder(binutil.HexDecoder)
	require.NoError(t, err, "could not create hex decoder")

	_, err = dec.DecodeString("68656c6c6f20776f726c64")
	require.NoError(t, err, "could not decode hex correctly")
}
