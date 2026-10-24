package freeD

import (
	"errors"
	"math"
)

type Position struct {
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

func (p Position) MarshalBinary() ([]byte, error) {
	bytes := p.Encode()
	return bytes[:], nil
}

func (p *Position) UnmarshalBinary(data []byte) error {
	if len(data) < 29 {
		return errors.New("data too short to be a valid FreeD position message")
	}

	if data[0] != 0xd1 {
		return errors.New("data is not a FreeD position message")
	}

	checksum := checksum(data[0 : len(data)-1])

	if checksum != data[len(data)-1] {
		return errors.New("FreeD packet checksum failed")
	}

	id := data[1]
	pan := freeDUnitsToRotation(data[2], data[3], data[4])
	tilt := freeDUnitsToRotation(data[5], data[6], data[7])
	roll := freeDUnitsToRotation(data[8], data[9], data[10])

	posX := freeDUnitsToPosition(data[11], data[12], data[13])
	posY := freeDUnitsToPosition(data[14], data[15], data[16])
	posZ := freeDUnitsToPosition(data[17], data[18], data[19])

	zoom := int32(data[20])*65536 + int32(data[21])*256 + int32(data[22])
	focus := int32(data[23])*65536 + int32(data[24])*256 + int32(data[25])

	p.ID = id
	p.Pan = pan
	p.Tilt = tilt
	p.Roll = roll
	p.PosX = posX
	p.PosY = posY
	p.PosZ = posZ
	p.Zoom = zoom
	p.Focus = focus

	return nil
}

func (p Position) Encode() [29]byte {
	var bytes [29]byte

	bytes[0] = 0xd1
	bytes[1] = p.ID

	copy(bytes[2:5], rotationToFreeDUnits(p.Pan))
	copy(bytes[5:8], rotationToFreeDUnits(p.Tilt))
	copy(bytes[8:11], rotationToFreeDUnits(p.Roll))

	copy(bytes[11:14], positionToFreeDUnits(p.PosX))
	copy(bytes[14:17], positionToFreeDUnits(p.PosY))
	copy(bytes[17:20], positionToFreeDUnits(p.PosZ))

	bytes[20] = uint8(math.Trunc(float64(p.Zoom / 65536)))
	bytes[21] = uint8(int64(p.Zoom/256) % 256)
	bytes[22] = uint8(int32(p.Zoom) % 256)

	bytes[23] = uint8(math.Trunc(float64(p.Focus / 65536)))
	bytes[24] = uint8(int64(p.Focus/256) % 256)
	bytes[25] = uint8(int32(p.Focus) % 256)

	// spare area?
	bytes[26] = 0
	bytes[27] = 0

	bytes[28] = checksum(bytes[0:28])
	return bytes
}

func Decode(bytes [29]byte) (*Position, error) {
	switch bytes[0] {
	case 0xd1:
		// FreeD position message
		position := &Position{}
		err := position.UnmarshalBinary(bytes[:])
		if err != nil {
			return nil, err
		}
		return position, nil
	default:
		return nil, errors.New("only FreeD position messages are currently supported")
	}
}
