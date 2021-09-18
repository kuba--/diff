package diff

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeltaAdd(t *testing.T) {
	require := require.New(t)

	const (
		strongSize = byte(4)
		blockSize  = uint32(11)

		oldText = `ala ma kota,kot ma ale,lal al ala,tyl e`
		newText = `toj es tto,ala ma kota,kot ma ale,lal al ala,tyl e`
	)
	delta := []*DeltaInstruction{
		{DeltaInstructionHeader: DeltaInstructionHeader{From: FromNew, Offset: 0, Size: uint64(blockSize)}},
		{DeltaInstructionHeader: DeltaInstructionHeader{From: FromOld, Offset: 0, Size: uint64(44)}},
	}

	oldReader := bytes.NewBufferString(oldText)
	oldBuffer := bytes.NewBuffer(nil)

	sig, err := WriteSignature(oldReader, oldBuffer, blockSize, strongSize)
	require.NoError(err)

	newReader := bytes.NewBufferString(newText)
	deltaBuffer := bytes.NewBuffer(nil)

	err = WriteDelta(sig, newReader, deltaBuffer)
	require.NoError(err)

	deltaReader := bytes.NewBuffer(deltaBuffer.Bytes())
	instr, err := ReadDelta(deltaReader)
	require.NoError(err)
	require.Len(instr, len(delta))

	for i, in := range instr {
		require.EqualValues(delta[i].DeltaInstructionHeader, in.DeltaInstructionHeader)
	}
}

func TestDeltaShift(t *testing.T) {
	require := require.New(t)

	const (
		strongSize = byte(8)
		blockSize  = uint32(11)

		oldText = `ala ma kotakot ma ale,lal al ala,tyl e`
		newText = `kot ma ale,ala ma kota,lal al ala,tyl e`
	)
	delta := []*DeltaInstruction{
		{DeltaInstructionHeader: DeltaInstructionHeader{From: FromOld, Offset: 11, Size: uint64(blockSize)}},
		{DeltaInstructionHeader: DeltaInstructionHeader{From: FromOld, Offset: 0, Size: uint64(blockSize)}},
		{DeltaInstructionHeader: DeltaInstructionHeader{From: FromNew, Offset: 0, Size: uint64(1)}},
		{DeltaInstructionHeader: DeltaInstructionHeader{From: FromOld, Offset: 22, Size: uint64(blockSize + blockSize)}},
	}

	oldReader := bytes.NewBufferString(oldText)
	oldBuffer := bytes.NewBuffer(nil)

	sig, err := WriteSignature(oldReader, oldBuffer, blockSize, strongSize)
	require.NoError(err)

	newReader := bytes.NewBufferString(newText)
	deltaBuffer := bytes.NewBuffer(nil)

	err = WriteDelta(sig, newReader, deltaBuffer)
	require.NoError(err)

	deltaReader := bytes.NewBuffer(deltaBuffer.Bytes())
	instr, err := ReadDelta(deltaReader)
	require.NoError(err)
	// require.Len(instr, len(delta))

	for i, in := range instr {
		t.Logf("[%d]: %+v\n", i, in.DeltaInstructionHeader)
		require.EqualValues(delta[i].DeltaInstructionHeader, in.DeltaInstructionHeader)
	}
}

func TestDeltaDelete(t *testing.T) {
	require := require.New(t)

	const (
		strongSize = byte(4)
		blockSize  = uint32(11)

		oldText = `ala ma kota,kot ma ale,lal al ala,tyl e`
		newText = `toj es tto,ala ma kota,tyl e`
	)
	delta := []*DeltaInstruction{
		{DeltaInstructionHeader: DeltaInstructionHeader{From: FromNew, Offset: 0, Size: uint64(blockSize)}},
		{DeltaInstructionHeader: DeltaInstructionHeader{From: FromOld, Offset: 0, Size: uint64(blockSize)}},
		{DeltaInstructionHeader: DeltaInstructionHeader{From: FromOld, Offset: 33, Size: uint64(blockSize)}},
	}

	oldReader := bytes.NewBufferString(oldText)
	oldBuffer := bytes.NewBuffer(nil)

	sig, err := WriteSignature(oldReader, oldBuffer, blockSize, strongSize)
	require.NoError(err)

	newReader := bytes.NewBufferString(newText)
	deltaBuffer := bytes.NewBuffer(nil)

	err = WriteDelta(sig, newReader, deltaBuffer)
	require.NoError(err)

	deltaReader := bytes.NewBuffer(deltaBuffer.Bytes())
	instr, err := ReadDelta(deltaReader)
	require.NoError(err)
	require.Len(instr, len(delta))

	for i, in := range instr {
		require.EqualValues(delta[i].DeltaInstructionHeader, in.DeltaInstructionHeader)
	}
}

func TestDeltaAddShiftDelete(t *testing.T) {
	require := require.New(t)

	const (
		strongSize = byte(4)
		blockSize  = uint32(11)

		oldText = `ala ma kota,1234567890,kot ma ale,lal al ala,tyl e`
		newText = `toj es tto,ala ma kota,1234567890,tyl e`
	)
	delta := []*DeltaInstruction{
		{DeltaInstructionHeader: DeltaInstructionHeader{From: FromNew, Offset: 0, Size: uint64(blockSize)}},
		{DeltaInstructionHeader: DeltaInstructionHeader{From: FromOld, Offset: 0, Size: uint64(blockSize + blockSize)}},
		{DeltaInstructionHeader: DeltaInstructionHeader{From: FromOld, Offset: 44, Size: uint64(blockSize)}},
	}

	oldReader := bytes.NewBufferString(oldText)
	oldBuffer := bytes.NewBuffer(nil)

	sig, err := WriteSignature(oldReader, oldBuffer, blockSize, strongSize)
	require.NoError(err)

	newReader := bytes.NewBufferString(newText)
	deltaBuffer := bytes.NewBuffer(nil)

	err = WriteDelta(sig, newReader, deltaBuffer)
	require.NoError(err)

	deltaReader := bytes.NewBuffer(deltaBuffer.Bytes())
	instr, err := ReadDelta(deltaReader)
	require.NoError(err)
	require.Len(instr, len(delta))

	for i, in := range instr {
		require.EqualValues(delta[i].DeltaInstructionHeader, in.DeltaInstructionHeader)
	}
}
