package freeD_test

import (
	"testing"

	freeD "github.com/jwetzell/free-d-go"
	"github.com/stretchr/testify/require"
)

func TestGoodPositionDecoding(t *testing.T) {
	testCases := []struct {
		description string
		bytes       [29]byte
		expected    freeD.Position
	}{
		{
			description: "simple position message",
			bytes: [29]byte{0xd1, 0x01, 0x5a, 0x00, 0x00, 0x2d, 0x00, 0x00, 0xa6, 0x00, 0x00, 0x7f, 0xff, 0x40, 0x7f, 0xff, 0x80, 0x7f, 0xff,
				0xc0, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x00, 0x00, 50,
			},
			expected: freeD.Position{
				ID:    1,
				Pan:   180,
				Tilt:  90,
				Roll:  -180,
				PosX:  131069,
				PosY:  131070,
				PosZ:  131071,
				Zoom:  66051,
				Focus: 263430,
			},
		},
	}

	for _, testCase := range testCases {

		actual := &freeD.Position{}

		err := actual.UnmarshalBinary(testCase.bytes[:])

		require.NoError(t, err)

		require.Equal(t, &testCase.expected, actual)
	}
}

func TestGoodPositionEncoding(t *testing.T) {
	testCases := []struct {
		description string
		expected    [29]byte
		message     freeD.Position
	}{
		{
			description: "simple position message",
			expected: [29]byte{0xd1, 0x01, 0x5a, 0x00, 0x00, 0x2d, 0x00, 0x00, 0xa6, 0x00, 0x00, 0x7f, 0xff, 0x40, 0x7f, 0xff, 0x80, 0x7f, 0xff,
				0xc0, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x00, 0x00, 50,
			},
			message: freeD.Position{
				ID:    1,
				Pan:   180,
				Tilt:  90,
				Roll:  -180,
				PosX:  131069,
				PosY:  131070,
				PosZ:  131071,
				Zoom:  66051,
				Focus: 263430,
			},
		},
	}

	for _, testCase := range testCases {

		actual := testCase.message.Encode()

		require.Equal(t, testCase.expected, actual)
	}
}

func BenchmarkMessageDecoding(b *testing.B) {
	data := [29]byte{0xd1, 0x01, 0x5a, 0x00, 0x00, 0x2d, 0x00, 0x00, 0xa6, 0x00, 0x00, 0x7f, 0xff, 0x40, 0x7f, 0xff, 0x80, 0x7f, 0xff,
		0xc0, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x00, 0x00, 50,
	}
	for b.Loop() {
		freeD.Decode(data)
	}
}

func BenchmarkMessageEncoding(b *testing.B) {
	message := freeD.Position{
		ID:    1,
		Pan:   180,
		Tilt:  90,
		Roll:  -180,
		PosX:  131069,
		PosY:  131070,
		PosZ:  131071,
		Zoom:  66051,
		Focus: 263430,
	}
	for b.Loop() {
		message.Encode()
	}
}

func FuzzMessageDecoding(f *testing.F) {
	f.Add([]byte{0xd1, 0x01, 0x5a, 0x00, 0x00, 0x2d, 0x00, 0x00, 0xa6, 0x00, 0x00, 0x7f, 0xff, 0x40, 0x7f, 0xff, 0x80, 0x7f, 0xff,
		0xc0, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x00, 0x00, 50,
	})
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) != 29 {
			return
		}
		var dataArray [29]byte
		copy(dataArray[:], data)
		freeD.Decode(dataArray)
	})
}
