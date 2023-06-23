package binutil_test

import (
	"testing"

	"github.com/bbengfort/binutil"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

func TestULID(t *testing.T) {
	testCases := []string{
		"00000000000000000000000000",
		"01H3MH1EP15QT769GDQFQ1E7T4",
		ulid.Make().String(),
	}

	for i, tc := range testCases {
		ds := &binutil.ULID{}
		es, err := ds.DecodeString(tc)
		require.NoError(t, err, "could not decode string for test case %d", i)

		data, err := es.EncodeBinary()
		require.NoError(t, err, "could not encode binary for test case %d", i)
		require.Len(t, data, 16, "expected 16 bytes of data for test case %d", i)

		db := &binutil.ULID{}
		eb, err := db.DecodeBinary(data)
		require.NoError(t, err, "could not decode binary for test case %d", i)

		s, err := eb.EncodeString()
		require.NoError(t, err, "could not encode string for test case %d", i)
		require.Equal(t, tc, s, "expected encoded string to match test case %d", i)
	}
}

func TestRegisteredULID(t *testing.T) {
	dec, err := binutil.NewDecoder(binutil.ULIDDecoder)
	require.NoError(t, err, "could not create ulid decoder")

	_, err = dec.DecodeString("01H3MH1EP15QT769GDQFQ1E7T4")
	require.NoError(t, err, "could not decode ulid correctly")
}
