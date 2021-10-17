package diff

const rollCharOffset = 31

// checksum32 was taken from librsync:
// https://github.com/librsync/librsync/blob/master/src/rollsum.c
func checksum32(p []byte) uint32 {
	s1 := uint16(0)
	s2 := uint16(0)
	l := len(p)
	for n := 0; n < l; {
		if n+15 < l {
			for i := 0; i < 16; i++ {
				s1 += uint16(p[n+i])
				s2 += s1
			}
			n += 16
		} else {
			s1 += uint16(p[n])
			s2 += s1
			n += 1
		}
	}

	s1 += uint16(l * rollCharOffset)
	s2 += uint16(((l * (l + 1)) / 2) * rollCharOffset)
	return (uint32(s2) << 16) | (uint32(s1) & 0xffff)
}

// rolling buffer is a circular buffer which calculates rolling checksum for written bytes.
// rolling checksum algorithm is heavily inspired by librsync:
// https://github.com/librsync/librsync/blob/master/src/rollsum.h
type rollBuffer struct {
	buf    []byte
	size   int
	pos    int
	count  int
	s1, s2 uint16
}

func newRollBuffer(size int) *rollBuffer {
	return &rollBuffer{
		buf:  make([]byte, size),
		pos:  0,
		size: size,
	}
}

func (rb *rollBuffer) reset() {
	rb.pos = 0
	rb.count = 0
	rb.s1 = 0
	rb.s2 = 0
}

// writeByte writes a new byte (in) to the buffer and returns (potentially) overwritten byte
func (rb *rollBuffer) writeByte(in byte) (out byte, overwrote bool) {
	overwrote = rb.count >= rb.size
	if overwrote {
		// rotate (in/out)
		out = rb.buf[rb.pos]
		rb.s1 += uint16(in) - uint16(out)
		rb.s2 += rb.s1 - uint16(rb.size)*(uint16(out)+uint16(rollCharOffset))
	} else {
		// rollin
		rb.s1 += uint16(in) + uint16(rollCharOffset)
		rb.s2 += rb.s1
	}

	rb.buf[rb.pos] = in
	rb.count++
	rb.pos = (rb.pos + 1) % rb.size
	return
}

func (rb *rollBuffer) bytes() []byte {
	if rb.count == 0 {
		return nil
	}

	if rb.count >= rb.size {
		if rb.pos == 0 {
			return rb.buf
		}
		buf := make([]byte, rb.size)
		copy(buf, rb.buf[rb.pos:])
		copy(buf[rb.size-rb.pos:], rb.buf[:rb.pos])
		return buf
	}

	return rb.buf[:rb.pos]
}

func (rb *rollBuffer) checksum32() uint32 {
	return (uint32(rb.s2) << 16) | (uint32(rb.s1) & 0xffff)
}
