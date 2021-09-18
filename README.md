# diff
A delta encoding library heavily inspired by [rdiff](https://linux.die.net/man/1/rdiff).
```
signature(basis-file) -> sig-file

delta(sig-file, new-file) -> delta-file

patch(basis-file, delta-file) -> recreated-file
```

The idea of _rolling checksum_ algorithm (_rollsum_) was taken from [librsync](https://github.com/librsync/librsync).

### API and File Format
- Signature
```go
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

diff.WriteSignature(basisReader io.Reader, signatureWriter io.Writer, blockSize uint32, strongSize byte) (*diff.Signature, error)
diff.ReadSignature(signatureReader io.Reader) (*diff.Signature, error)

func (sig *Signature) Lookup(weak uint32) (strong []byte, offset uint64, blockSize uint32, ok bool)
```

```
// header
{block size: 4 bytes, strong checksum size: 1 byte}
// checksum
{weak checksum: 4 bytes, strong checksum: StrongSize bytes}
{weak checksum: 4 bytes, strong checksum: StrongSize bytes}
...
{weak checksum: 4 bytes, strong checksum: StrongSize bytes}
```

- Delta
```go
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

diff.WriteDelta(signature *diff.Signature, newReader io.Reader, deltaWriter io.Writer) error
diff.ReadDelta(r io.Reader) (delta diff.Delta, err error)
diff.ReadDeltaInstructionHeader(r io.Reader) (header diff.DeltaInstructionHeader, err error)
```

```
// instruction
{from: 1 byte, offset: 8 bytes, size: 8 bytes}
// data
...

// instruction
{from: 1 byte, offset: 8 bytes, size: 8 bytes}
// data
...
```

- Patch
```go
diff.Patch(basisReaderSeeker io.ReadSeeker, deltaReader io.Reader, newWriter io.Writer) error
```

### Usage
```
go build ./cmd/signature
./signature [-b block size] old-file signature-file

go build ./cmd/delta
./delta signature-file new-file delta-file

go build ./cmd/patch
./patch old-file delta-file new-file
```