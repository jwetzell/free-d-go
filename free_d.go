package freeD

import (
	"bytes"
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

func Encode(message FreeDPosition) []byte {
	bytes := []byte{}

	bytes = append(bytes, 0xd1)
	bytes = append(bytes, message.ID)

	bytes = append(bytes, rotationToFreeDUnits(message.Pan)...)
	bytes = append(bytes, rotationToFreeDUnits(message.Tilt)...)
	bytes = append(bytes, rotationToFreeDUnits(message.Roll)...)

	bytes = append(bytes, positionToFreeDUnits(message.PosX)...)
	bytes = append(bytes, positionToFreeDUnits(message.PosY)...)
	bytes = append(bytes, positionToFreeDUnits(message.PosZ)...)

	bytes = append(bytes, uint8(math.Trunc(float64(message.Zoom/65536))))
	bytes = append(bytes, uint8(int64(message.Zoom/256)%256))
	bytes = append(bytes, uint8(int32(message.Zoom)%256))

	bytes = append(bytes, uint8(math.Trunc(float64(message.Focus/65536))))
	bytes = append(bytes, uint8(int64(message.Focus/256)%256))
	bytes = append(bytes, uint8(int32(message.Focus)%256))

	// spare area?
	bytes = append(bytes, 0, 0)

	bytes = append(bytes, checksum(bytes[0:28]))
	return bytes
}

func Decode(bytes []byte) (FreeDPosition, error) {
	if len(bytes) != 29 {
		return FreeDPosition{}, errors.New("FreeD packet must be 29 bytes long")
	}

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

	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, units)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()[0:3]
}

func positionToFreeDUnits(position float32) []byte {
	units := int32(position * 64 * 256)

	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, units)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()[0:3]
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
