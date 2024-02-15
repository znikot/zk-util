package misc

import (
	"encoding/binary"
)

// convert int to byte slice
func IntToBytes(i int) []byte {
	return []byte{
		byte(i >> 24),
		byte(i >> 16),
		byte(i >> 8),
		byte(i),
	}
}

// convert byte slice to int
func BytesToInt(b []byte) int {
	return int(int32(b[3]) | int32(b[2])<<8 | int32(b[1])<<16 | int32(b[0])<<24)
}

// convert int64 to byte slice
func Int64ToBytes(i int64) []byte {
	return []byte{
		byte(i >> 56),
		byte(i >> 48),
		byte(i >> 40),
		byte(i >> 32),
		byte(i >> 24),
		byte(i >> 16),
		byte(i >> 8),
		byte(i),
	}
}

// convert byte slice to int64
func BytesToInt64(b []byte) int64 {
	return int64(b[7]) | int64(b[6])<<8 | int64(b[5])<<16 | int64(b[4])<<24 | int64(b[3])<<32 | int64(b[2])<<40 | int64(b[1])<<48 | int64(b[0])<<56
}

// convert int16 to byte slice
func Int16ToBytes(i int16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint16(i))
	return bytes
}

// convert byte slice to int16
func BytesToInt16(b []byte) int16 {
	return int16(binary.BigEndian.Uint16(b))
}
