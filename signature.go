package diff

import (
	"errors"
	"io"
)

type (
	Signature struct {
		signatureHeader
		signatureChecksum
	}

	signatureHeader struct {
		BlockSize  uint32
		StrongSize byte
	}

	signatureChecksum struct {
		weak   map[uint32]int
		strong [][]byte
	}
)

// WriteSignature generates the signature of a basis reader, and writes it out to signatureWriter.
func WriteSignature(basisReader io.Reader, signatureWriter io.Writer, blockSize uint32, strongSize byte) (*Signature, error) {
	if blockSize == 0 {
		return nil, errors.New("block size must be > 0")
	}
	if strongSize == 0 {
		return nil, errors.New("hash size must be > 0")
	}

	header, err := writeSignatureHeader(signatureWriter, blockSize, strongSize)
	if err != nil {
		return nil, err
	}
	checksum, err := writeSignatureChecksum(basisReader, signatureWriter, blockSize, strongSize)
	if err != nil {
		return nil, err
	}
	return &Signature{header, checksum}, nil
}

// ReadSignature reads the signature from signatureReader.
func ReadSignature(signatureReader io.Reader) (*Signature, error) {
	header, err := readSignatureHeader(signatureReader)
	if err != nil {
		return nil, err
	}

	checksum, err := readSignatureChecksum(signatureReader, header.StrongSize)
	if err != nil {
		return nil, err
	}

	return &Signature{header, checksum}, nil
}

// Lookup retrieves to block for a given weak checksum.
func (sig *Signature) Lookup(weak uint32) (strong []byte, offset uint64, blockSize uint32, ok bool) {
	var idx int
	idx, ok = sig.weak[weak]
	if !ok {
		return
	}

	strong = sig.strong[idx]
	offset = uint64(idx) * uint64(sig.BlockSize)
	blockSize = sig.BlockSize
	return
}

func writeSignatureHeader(w io.Writer, blockSize uint32, strongSize byte) (header signatureHeader, err error) {
	var b [4 + 1]byte
	// block size
	ByteOrder.PutUint32(b[:4], blockSize)
	// strong size
	b[4] = strongSize

	if _, err = w.Write(b[:]); err != nil {
		return
	}
	header = signatureHeader{BlockSize: blockSize, StrongSize: strongSize}
	return
}

func readSignatureHeader(r io.Reader) (header signatureHeader, err error) {
	var b [4 + 1]byte
	if _, err = r.Read(b[:]); err != nil {
		return
	}

	// block size
	header.BlockSize = ByteOrder.Uint32(b[:4])
	// strong size
	header.StrongSize = b[4]
	return
}

func writeSignatureChecksum(r io.Reader, w io.Writer, blockSize uint32, strongSize byte) (signatureChecksum, error) {
	checksum := signatureChecksum{weak: make(map[uint32]int)}

	var weak [4]byte
	buf := make([]byte, blockSize)
	h := NewHash()
	for i := 0; ; i++ {
		n, err := io.ReadFull(r, buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			if err != io.ErrUnexpectedEOF {
				return signatureChecksum{}, err
			}
		}

		// write weak checksum
		v := checksum32(buf[:n])
		ByteOrder.PutUint32(weak[:], v)
		if _, err = w.Write(weak[:]); err != nil {
			return signatureChecksum{}, err
		}
		checksum.weak[v] = i

		// write strong checksum
		h.Reset()
		if _, err = h.Write(buf[:n]); err != nil {
			return signatureChecksum{}, err
		}
		strong := h.Sum(nil)[:strongSize]
		if _, err = w.Write(strong); err != nil {
			return signatureChecksum{}, err
		}
		checksum.strong = append(checksum.strong, make([]byte, strongSize))
		copy(checksum.strong[i], strong)
	}

	return checksum, nil
}

func readSignatureChecksum(r io.Reader, strongSize byte) (signatureChecksum, error) {
	checksum := signatureChecksum{weak: make(map[uint32]int)}

	var weak [4]byte
	strong := make([]byte, strongSize)
	for i := 0; ; i++ {
		// read weak checksum
		if _, err := r.Read(weak[:]); err != nil {
			if err == io.EOF {
				break
			}
			return signatureChecksum{}, err
		}
		// read strong checksum
		if _, err := r.Read(strong); err != nil {
			if err == io.EOF {
				break
			}
			return signatureChecksum{}, err
		}

		checksum.weak[ByteOrder.Uint32(weak[:])] = i
		checksum.strong = append(checksum.strong, make([]byte, strongSize))
		copy(checksum.strong[i], strong)
	}

	return checksum, nil
}
