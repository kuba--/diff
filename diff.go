package diff

import (
	"crypto/md5"
	"encoding/binary"
)

var (
	ByteOrder = binary.BigEndian
	NewHash   = md5.New
)
