package binutil_test

import (
	"testing"

	"github.com/bbengfort/binutil"
	"github.com/stretchr/testify/require"
)

func TestBase64(t *testing.T) {
	makeTestForScheme := func(scheme binutil.Base64Scheme) func(t *testing.T) {
		return func(t *testing.T) {
			// Test scheme string
			require.NotEmpty(t, scheme.String(), "expected scheme string to be returned")
			require.NotEqual(t, "unknown", scheme.String(), "expected scheme string to not be unknown")

			// Create encoder and fixtures for tests
			b64 := binutil.NewBase64(scheme)
			fixtures := fixtures()

			for _, fixture := range fixtures {
				eb, err := b64.DecodeBinary(fixture.data)
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

				es, err := b64.DecodeString(s)
				require.NoError(t, err, "could not decode string for fixture %q", fixture.name)

				data, err = es.EncodeBinary()
				require.NoError(t, err, "could not encode binary from string for fixture %q", fixture.name)
				require.Equal(t, fixture.data, data, "expected unchanged binary data from decoded string for fixture %q", fixture.name)
			}
		}
	}

	t.Run("Std", makeTestForScheme(binutil.B64SchemeStd))
	t.Run("RawStd", makeTestForScheme(binutil.B64SchemeRawStd))
	t.Run("URL", makeTestForScheme(binutil.B64SchemeURL))
	t.Run("RawURL", makeTestForScheme(binutil.B64SchemeRawURL))
}
