package freeD

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMessageDecoding(t *testing.T) {
	testCases := []struct {
		description string
		bytes       []byte
		expected    FreeDPosition
	}{
		{
			description: "simple position message",
			bytes: []byte{0xd1, 0x01, 0x5a, 0x00, 0x00, 0x2d, 0x00, 0x00, 0xa6, 0x00, 0x00, 0x7f, 0xff, 0x40, 0x7f, 0xff, 0x80, 0x7f, 0xff,
				0xc0, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x00, 0x00, 50,
			},
			expected: FreeDPosition{
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

		actual, err := Decode(testCase.bytes)

		if err != nil {
			t.Errorf("Test '%s' failed to decode chunk properly", testCase.description)
			fmt.Println(err)
		}

		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("Test '%s' failed to decode chunk properly", testCase.description)
			fmt.Printf("expected: %+v\n", testCase.expected)
			fmt.Printf("actual: %+v\n", actual)
		}
	}
}

func TestMessageEncoding(t *testing.T) {
	testCases := []struct {
		description string
		expected    []byte
		message     FreeDPosition
	}{
		{
			description: "simple position message",
			expected: []byte{0xd1, 0x01, 0x5a, 0x00, 0x00, 0x2d, 0x00, 0x00, 0xa6, 0x00, 0x00, 0x7f, 0xff, 0x40, 0x7f, 0xff, 0x80, 0x7f, 0xff,
				0xc0, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x00, 0x00, 50,
			},
			message: FreeDPosition{
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

		actual := Encode(testCase.message)

		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("Test '%s' failed to decode chunk properly", testCase.description)
			fmt.Printf("expected: %+v\n", testCase.expected)
			fmt.Printf("actual: %+v\n", actual)
		}
	}
}
