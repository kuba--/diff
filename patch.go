package diff

import (
	"io"
)

func Patch(basisReaderSeeker io.ReadSeeker, deltaReader io.Reader, newWriter io.Writer) error {
	for {
		i, err := ReadDeltaInstructionHeader(deltaReader)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if i.From == FromOld {
			if _, err = basisReaderSeeker.Seek(int64(i.Offset), io.SeekStart); err != nil {
				return err
			}
			if _, err = io.CopyN(newWriter, basisReaderSeeker, int64(i.Size)); err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
		} else if i.From == FromNew {
			if _, err = io.CopyN(newWriter, deltaReader, int64(i.Size)); err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
		}
	}

	return nil
}
