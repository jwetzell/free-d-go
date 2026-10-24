package freeD

import "encoding/binary"

func checksum(bytes []byte) uint8 {
	var checksum uint8 = 0x40

	for _, value := range bytes {
		checksum -= value
	}
	return checksum
}

func rotationToFreeDUnits(rotation float32) []byte {
	units := int32(rotation * 32768 * 256)

	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(units))
	return buf[0:3]
}

func positionToFreeDUnits(position float32) []byte {
	units := int32(position * 64 * 256)

	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(units))
	return buf[0:3]
}

func freeDUnitsToRotation(upper uint8, middle uint8, lower uint8) float32 {
	var units int32
	units |= int32(lower) << 8
	units |= int32(middle) << 16
	units |= int32(upper) << 24

	return float32(units) / 256 / 32768
}

func freeDUnitsToPosition(upper uint8, middle uint8, lower uint8) float32 {
	var units int32
	units |= int32(lower) << 8
	units |= int32(middle) << 16
	units |= int32(upper) << 24

	return float32(units) / 256 / 64
}
