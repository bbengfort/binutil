package binutil_test

import (
	"testing"

	"github.com/bbengfort/binutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUUID(t *testing.T) {
	testCases := []string{
		"00000000-0000-0000-0000-000000000000",
		"36a15b36-3e89-45cc-ae97-4813ce4ead77",
		uuid.New().String(),
	}

	for i, tc := range testCases {
		ds := &binutil.UUID{}
		es, err := ds.DecodeString(tc)
		require.NoError(t, err, "could not decode string for test case %d", i)

		data, err := es.EncodeBinary()
		require.NoError(t, err, "could not encode binary for test case %d", i)
		require.Len(t, data, 16, "expected 16 bytes of data for test case %d", i)

		db := &binutil.UUID{}
		eb, err := db.DecodeBinary(data)
		require.NoError(t, err, "could not decode binary for test case %d", i)

		s, err := eb.EncodeString()
		require.NoError(t, err, "could not encode string for test case %d", i)
		require.Equal(t, tc, s, "expected encoded string to match test case %d", i)
	}
}
