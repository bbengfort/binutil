package binutil_test

import (
	"bytes"
	crand "crypto/rand"
	"testing"

	"github.com/bbengfort/binutil"
	"github.com/stretchr/testify/require"
)

func TestPipelineStr2Str(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		steps    []any
	}{
		{"01H3W1T4BNATG1KGP7S817K4BF", "01H3W1T4BNATG1KGP7S817K4BF", []any{"ulid"}},
		{"01H3W1T4BNATG1KGP7S817K4BF", "AYj4HRF1VqAZwsfKAnmRbw==", []any{"ulid", "b64"}},
		{"AYj4HRF1VqAZwsfKAnmRbw==", "01H3W1T4BNATG1KGP7S817K4BF", []any{"b64", "ulid"}},
		{"01H3W1T4BNATG1KGP7S817K4BF", "0188f81d117556a019c2c7ca0279916f", []any{"ulid", "hex"}},
		{"0188f81d117556a019c2c7ca0279916f", "01H3W1T4BNATG1KGP7S817K4BF", []any{"hex", "ulid"}},
		{"a1372ed62623e0037e46c31535a407041e48c21cb240acf5bc8863", "oTcu1iYj4AN+RsMVNaQHBB5IwhyyQKz1vIhj", []any{"hex", "b64"}},
		{"oTcu1iYj4AN+RsMVNaQHBB5IwhyyQKz1vIhj", "a1372ed62623e0037e46c31535a407041e48c21cb240acf5bc8863", []any{"b64", "hex"}},
		{"3ecb2f46-0242-4642-bdef-91d191650369", "PssvRgJCRkK975HRkWUDaQ==", []any{"uuid", "b64"}},
		{"PssvRgJCRkK975HRkWUDaQ==", "3ecb2f46-0242-4642-bdef-91d191650369", []any{"b64", "uuid"}},
		{"3ecb2f46-0242-4642-bdef-91d191650369", "3ecb2f4602424642bdef91d191650369", []any{"uuid", "hex"}},
		{"3ecb2f4602424642bdef91d191650369", "3ecb2f46-0242-4642-bdef-91d191650369", []any{"hex", "uuid"}},
		{"01H3W1T4BNATG1KGP7S817K4BF", "0188f81d-1175-56a0-19c2-c7ca0279916f", []any{"ulid", "uuid"}},
		{"0188f81d-1175-56a0-19c2-c7ca0279916f", "01H3W1T4BNATG1KGP7S817K4BF", []any{"uuid", "ulid"}},
		{"0188f81d-1175-56a0-19c2-c7ca0279916f", "0188f81d-1175-56a0-19c2-c7ca0279916f", []any{"uuid", "hex", "b64", "uuid"}},
	}

	for i, tc := range testCases {
		pipe, err := binutil.New(tc.steps...)
		require.NoError(t, err, "could not make pipeline for test case %d", i)

		actual, err := pipe.Str2Str(tc.input)
		require.NoError(t, err, "could not convert str to str for test case %d", i)
		require.Equal(t, tc.expected, actual, "incorrect str2str conversion for test case %d", i)
	}
}

func TestPipelineBin2Bin(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected []byte
		steps    []any
	}{
		{[]byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, []byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, []any{"uuid"}},
		{[]byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, []byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, []any{"hex"}},
		{[]byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, []byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, []any{"ulid", "uuid"}},
		{[]byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, []byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, []any{"hex", "b64"}},
		{[]byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, []byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, []any{"ulid", "hex", "b64", "uuid"}},
	}

	for i, tc := range testCases {
		pipe, err := binutil.New(tc.steps...)
		require.NoError(t, err, "could not make pipeline for test case %d", i)

		actual, err := pipe.Bin2Bin(tc.input)
		require.NoError(t, err, "could not convert str to str for test case %d", i)
		require.Equal(t, tc.expected, actual, "incorrect str2str conversion for test case %d", i)
	}
}

func TestPipelineBin2Str(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected string
		steps    []any
	}{
		{[]byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, "54011b6f-eb92-84f6-0217-2da7bea9375a", []any{"uuid"}},
		{[]byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, "2M04DPZTWJGKV045SDMYZAJDTT", []any{"ulid"}},
		{[]byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, "VAEbb+uShPYCFy2nvqk3Wg", []any{"hex", "b64raw"}},
		{[]byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, "54011b6feb9284f602172da7bea9375a", []any{"uuid", "hex"}},
		{[]byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, "54011b6f-eb92-84f6-0217-2da7bea9375a", []any{"b64", "hex", "uuid"}},
	}

	for i, tc := range testCases {
		pipe, err := binutil.New(tc.steps...)
		require.NoError(t, err, "could not make pipeline for test case %d", i)

		actual, err := pipe.Bin2Str(tc.input)
		require.NoError(t, err, "could not convert str to str for test case %d", i)
		require.Equal(t, tc.expected, actual, "incorrect str2str conversion for test case %d", i)
	}
}

func TestPipelineStr2Bin(t *testing.T) {
	testCases := []struct {
		input    string
		expected []byte
		steps    []any
	}{
		{"54011b6f-eb92-84f6-0217-2da7bea9375a", []byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, []any{"uuid"}},
		{"54011b6f-eb92-84f6-0217-2da7bea9375a", []byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, []any{"uuid", "ulid"}},
		{"54011b6f-eb92-84f6-0217-2da7bea9375a", []byte{84, 1, 27, 111, 235, 146, 132, 246, 2, 23, 45, 167, 190, 169, 55, 90}, []any{"uuid", "hex", "b64", "ulid"}},
	}

	for i, tc := range testCases {
		pipe, err := binutil.New(tc.steps...)
		require.NoError(t, err, "could not make pipeline for test case %d", i)

		actual, err := pipe.Str2Bin(tc.input)
		require.NoError(t, err, "could not convert str to str for test case %d", i)
		require.Equal(t, tc.expected, actual, "incorrect str2str conversion for test case %d", i)
	}
}

func TestUnknownDecoder(t *testing.T) {
	dec, err := binutil.NewDecoder(" UnknownDECODER ")
	require.EqualError(t, err, "no registered decoder with the name \"unknowndecoder\"")
	require.Nil(t, dec, "expected returned decoder to be nil")
}

const (
	smallSize         = 16
	abnormalSmallSize = 23
	largeSize         = 4096
	abnormalLargeSize = 5167
)

type fixture struct {
	name string
	data []byte
}

// Create a list of test fixtures and test cases
func fixtures() []*fixture {
	tests := make([]*fixture, 0, 15)

	// Add some well structured fixtures
	tests = append(tests, &fixture{"countdown", []byte{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}})
	tests = append(tests, &fixture{"primes", []byte{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223, 227, 229, 233, 239, 241, 251}})

	// Add some zero valued fixtures
	tests = append(tests, &fixture{"empty, non-nil data", zeros(0)})
	tests = append(tests, &fixture{"small zeros", zeros(smallSize)})
	tests = append(tests, &fixture{"abnormal small zeros", zeros(abnormalSmallSize)})
	tests = append(tests, &fixture{"large zeros", zeros(largeSize)})
	tests = append(tests, &fixture{"abnormal large zeros", zeros(abnormalLargeSize)})

	// Lucky valued fixtures
	tests = append(tests, &fixture{"small lucky", repeat([]byte{8}, smallSize)})
	tests = append(tests, &fixture{"abnormal small lucky", repeat([]byte{8}, abnormalSmallSize)})
	tests = append(tests, &fixture{"large lucky", repeat([]byte{8}, largeSize)})
	tests = append(tests, &fixture{"abnormal large lucky", repeat([]byte{8}, abnormalLargeSize)})

	// Random data fixtures
	tests = append(tests, &fixture{"small rand", rand(smallSize)})
	tests = append(tests, &fixture{"abnormal small rand", rand(abnormalSmallSize)})
	tests = append(tests, &fixture{"large rand", rand(largeSize)})
	tests = append(tests, &fixture{"abnormal large rand", rand(abnormalLargeSize)})

	return tests
}

// Test fixture with all zeros of the specified length
func zeros(n int) []byte {
	return make([]byte, n)
}

// Test fixture with the same bytes of the specified length
func repeat(b []byte, n int) []byte {
	return bytes.Repeat(b, n)
}

// Random binary data of the specified length
func rand(n int) (out []byte) {
	out = make([]byte, n)
	crand.Read(out)
	return out
}
