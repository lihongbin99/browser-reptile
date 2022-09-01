package utils

func ToUint16(data []byte) uint16 {
	return uint16(data[0])<<8 +
		uint16(data[1])
}
