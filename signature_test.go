package diff

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignatureHeader(t *testing.T) {
	require := require.New(t)

	buf := bytes.NewBuffer(nil)
	header := signatureHeader{
		StrongSize: 8,
		BlockSize:  4096,
	}

	h1, err := writeSignatureHeader(buf, header.BlockSize, header.StrongSize)
	require.NoError(err)
	require.Equal(header, h1)

	h2, err := readSignatureHeader(buf)
	require.NoError(err)
	require.Equal(header, h2)
}

func TestSignatureChecksum(t *testing.T) {
	require := require.New(t)

	const (
		strongSize = 4
		blockSize  = 4

		text = `
	Lorem ipsum dolor sit amet, consectetur adipiscing elit,
	sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
	Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
	Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
	Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
	`
	)
	r := bytes.NewBufferString(text)
	rw := bytes.NewBuffer(nil)

	ch1, err := writeSignatureChecksum(r, rw, blockSize, strongSize)
	require.NoError(err)

	ch2, err := readSignatureChecksum(rw, strongSize)
	require.NoError(err)

	require.EqualValues(ch1, ch2)
}

func TestSignature(t *testing.T) {
	require := require.New(t)

	const (
		strongSize = 4
		blockSize  = 4

		text = `
	Lorem ipsum dolor sit amet, consectetur adipiscing elit,
	sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
	Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
	Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
	Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
	`
	)
	r := bytes.NewBufferString(text)
	rw := bytes.NewBuffer(nil)

	sig1, err := WriteSignature(r, rw, blockSize, strongSize)
	require.NoError(err)

	sig2, err := ReadSignature(rw)
	require.NoError(err)

	require.EqualValues(sig1, sig2)

	for weak := range sig1.weak {
		strong, _, size, ok := sig2.Lookup(weak)
		require.True(ok)
		require.Truef(size > 0, "size > 0")
		require.Len(strong, strongSize)
	}
}
