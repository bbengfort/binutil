package binutil_test

import (
	"bytes"
	crand "crypto/rand"
)

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
