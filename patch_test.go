package diff

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	strongSize = byte(4)
	blockSize  = uint32(11)

	basisText = `ala ma kota,1234567890,kot ma ale,lal al ala,tyl e`
	newText   = `toj es tto,ala ma kota,1234567890,tyl e`
)

var (
	basisFile *os.File
	newFile   *os.File
	sigFile   *os.File
	deltaFile *os.File
)

func setup(t *testing.T) {
	t.Helper()
	require := require.New(t)

	var err error

	basisFile, err = os.CreateTemp("", "basis")
	require.NoError(err)
	basisFile.WriteString(basisText)
	require.NoError(basisFile.Close())
	basisFile, err = os.Open(basisFile.Name())
	require.NoError(err)

	newFile, err = os.CreateTemp("", "new")
	require.NoError(err)
	newFile.WriteString(newText)
	require.NoError(newFile.Close())
	newFile, err = os.Open(newFile.Name())
	require.NoError(err)

	sigFile, err = os.CreateTemp("", "sig")
	require.NoError(err)
	sig, err := WriteSignature(basisFile, sigFile, blockSize, strongSize)
	require.NoError(err)

	deltaFile, err = os.CreateTemp("", "patch")
	require.NoError(err)

	err = WriteDelta(sig, newFile, deltaFile)
	require.NoError(err)
	require.NoError(deltaFile.Close())
	deltaFile, err = os.Open(deltaFile.Name())
	require.NoError(err)
}

func tearDown(t *testing.T) {
	t.Helper()

	basisFile.Close()
	os.Remove(basisFile.Name())

	newFile.Close()
	os.Remove(newFile.Name())

	sigFile.Close()
	os.Remove(sigFile.Name())

	deltaFile.Close()
	os.Remove(deltaFile.Name())
}

func TestPatch(t *testing.T) {
	require := require.New(t)
	setup(t)
	defer tearDown(t)

	buf := bytes.NewBuffer(nil)
	err := Patch(basisFile, deltaFile, buf)
	require.NoError(err)

	require.EqualValues(newText, buf.String())
}
