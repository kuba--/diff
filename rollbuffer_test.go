package diff

import "testing"

func TestChecksum32(t *testing.T) {
	str := "1234567890abcdefghijk"
	bstr := []byte(str)

	buf := newRollBuffer(4)
	for _, b := range bstr {
		buf.writeByte(b)
		if buf.count < buf.size {
			continue
		}

		c1 := checksum32(bstr[buf.count-buf.size : buf.count])
		c2 := buf.checksum32()
		t.Logf("checksum1(%s): %d, checksum2(%s): %d", bstr[buf.count-buf.size:buf.count], c1, buf.bytes(), c2)
		if c1 != c2 {
			t.Fatalf("expected: %d, got: %d", c1, c2)
		}
	}
}
