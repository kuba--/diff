package diff

import (
	"bufio"
	"bytes"
	"io"
)

const (
	FromOld = byte(0x0)
	FromNew = byte(0x1)
)

type (
	Delta = []*DeltaInstruction

	DeltaInstructionHeader struct {
		From   byte
		Offset uint64
		Size   uint64
	}

	DeltaInstruction struct {
		DeltaInstructionHeader
		Data []byte
	}
)

func WriteDelta(signature *Signature, newReader io.Reader, deltaWriter io.Writer) error {
	rd := bufio.NewReaderSize(newReader, int(signature.BlockSize))
	buf := newRollBuffer(int(signature.BlockSize))
	h := NewHash()

	i := &DeltaInstruction{}
	for {
		in, err := rd.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		out, overwrote := buf.writeByte(in)
		if buf.count < buf.size {
			continue
		}
		if overwrote {
			i.append(deltaWriter, &DeltaInstruction{
				DeltaInstructionHeader: DeltaInstructionHeader{From: FromNew, Size: uint64(1)},
				Data:                   []byte{out},
			})
		}

		weak := buf.checksum32()
		strong, offset, blocksize, ok := signature.Lookup(weak)
		// fmt.Printf("LOOKUP(%s):%d => [%v]: Offset:%d, Size:%d\n", buf.bytes(), weak, ok, offset, blocksize)

		if ok {
			// from old
			h.Reset()
			h.Write(buf.bytes())
			if bytes.Equal(strong, h.Sum(nil)[:signature.StrongSize]) {
				i.append(deltaWriter, &DeltaInstruction{
					DeltaInstructionHeader: DeltaInstructionHeader{From: FromOld,
						Offset: offset,
						Size:   uint64(blocksize),
					},
					Data: []byte{},
				})
			}
			buf.reset()
		}
	}

	// handle potential leftovers
	weak := buf.checksum32()
	strong, offset, blocksize, ok := signature.Lookup(weak)
	if ok {
		// from old
		h.Reset()
		h.Write(buf.bytes())
		if bytes.Equal(strong, h.Sum(nil)[:signature.StrongSize]) {
			i.append(deltaWriter, &DeltaInstruction{
				DeltaInstructionHeader: DeltaInstructionHeader{From: FromOld,
					Offset: offset,
					Size:   uint64(blocksize),
				},
				Data: []byte{},
			})
		}
		buf.reset()
	}

	for _, b := range buf.bytes() {
		i.append(deltaWriter, &DeltaInstruction{
			DeltaInstructionHeader: DeltaInstructionHeader{From: FromNew, Size: uint64(1)},
			Data:                   []byte{b},
		})
	}
	i.writeTo(deltaWriter)

	return nil
}

func ReadDelta(r io.Reader) (delta Delta, err error) {
	for {
		var i DeltaInstruction
		i.DeltaInstructionHeader, err = ReadDeltaInstructionHeader(r)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if i.From == FromNew && i.Size > 0 {
			i.Data = make([]byte, i.Size)
			if _, err = r.Read(i.Data); err != nil {
				return nil, err
			}
		}
		delta = append(delta, &i)
	}

	return delta, nil
}

func ReadDeltaInstructionHeader(r io.Reader) (header DeltaInstructionHeader, err error) {
	var b [1 + 8 + 8]byte
	if _, err = r.Read(b[:]); err != nil {
		return
	}

	header.From = b[0]
	header.Offset = ByteOrder.Uint64(b[1:9])
	header.Size = ByteOrder.Uint64(b[9:])
	return
}

func (i *DeltaInstruction) append(w io.Writer, next *DeltaInstruction) error {
	if next == nil || next.Size == 0 {
		return nil
	}

	if i.From != next.From {
		if err := i.writeTo(w); err != nil {
			return err
		}

		i.From = next.From
		i.Offset = next.Offset
		i.Size = next.Size
		i.Data = next.Data
		return nil
	}

	if i.From == FromNew {
		i.Data = append(i.Data, next.Data...)
		i.Size++
	} else if i.From == FromOld {
		if i.Offset+i.Size == next.Offset {
			// merge blocks
			i.Size += next.Size
		} else {
			if err := i.writeTo(w); err != nil {
				return err
			}

			i.Offset = next.Offset
			i.Size = next.Size
		}
	}

	return nil
}

func (i *DeltaInstruction) writeTo(w io.Writer) error {
	if i.Size == 0 {
		return nil
	}

	var b [1 + 8 + 8]byte
	b[0] = i.From
	ByteOrder.PutUint64(b[1:9], i.Offset)
	ByteOrder.PutUint64(b[9:], i.Size)
	if _, err := w.Write(b[:]); err != nil {
		return err
	}

	if i.From == FromNew && i.Data != nil {
		if _, err := w.Write(i.Data); err != nil {
			return err
		}
	}
	return nil
}
