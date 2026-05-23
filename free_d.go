package freeD

import (
	"encoding/binary"
	"errors"
	"math"
)

type FreeDPosition struct {
	ID    uint8
	Pan   float32
	Tilt  float32
	Roll  float32
	PosX  float32
	PosY  float32
	PosZ  float32
	Zoom  int32
	Focus int32
}

func Encode(message FreeDPosition) [29]byte {
	var bytes [29]byte

	bytes[0] = 0xd1
	bytes[1] = message.ID

	copy(bytes[2:5], rotationToFreeDUnits(message.Pan))
	copy(bytes[5:8], rotationToFreeDUnits(message.Tilt))
	copy(bytes[8:11], rotationToFreeDUnits(message.Roll))

	copy(bytes[11:14], positionToFreeDUnits(message.PosX))
	copy(bytes[14:17], positionToFreeDUnits(message.PosY))
	copy(bytes[17:20], positionToFreeDUnits(message.PosZ))

	bytes[20] = uint8(math.Trunc(float64(message.Zoom / 65536)))
	bytes[21] = uint8(int64(message.Zoom/256) % 256)
	bytes[22] = uint8(int32(message.Zoom) % 256)

	bytes[23] = uint8(math.Trunc(float64(message.Focus / 65536)))
	bytes[24] = uint8(int64(message.Focus/256) % 256)
	bytes[25] = uint8(int32(message.Focus) % 256)

	// spare area?
	bytes[26] = 0
	bytes[27] = 0

	bytes[28] = checksum(bytes[0:28])
	return bytes
}

func Decode(bytes [29]byte) (FreeDPosition, error) {
	if bytes[0] != 0xd1 {
		return FreeDPosition{}, errors.New("only FreeD position messages are currently supported")
	}

	checksum := checksum(bytes[0 : len(bytes)-1])

	if checksum != bytes[len(bytes)-1] {
		return FreeDPosition{}, errors.New("FreeD packet checksum failed")
	}

	id := bytes[1]
	pan := freeDUnitsToRotation(bytes[2], bytes[3], bytes[4])
	tilt := freeDUnitsToRotation(bytes[5], bytes[6], bytes[7])
	roll := freeDUnitsToRotation(bytes[8], bytes[9], bytes[10])

	posX := freeDUnitsToPosition(bytes[11], bytes[12], bytes[13])
	posY := freeDUnitsToPosition(bytes[14], bytes[15], bytes[16])
	posZ := freeDUnitsToPosition(bytes[17], bytes[18], bytes[19])

	zoom := int32(bytes[20])*65536 + int32(bytes[21])*256 + int32(bytes[22])
	focus := int32(bytes[23])*65536 + int32(bytes[24])*256 + int32(bytes[25])

	return FreeDPosition{
		ID:    id,
		Pan:   pan,
		Tilt:  tilt,
		Roll:  roll,
		PosX:  posX,
		PosY:  posY,
		PosZ:  posZ,
		Zoom:  zoom,
		Focus: focus,
	}, nil
}

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
