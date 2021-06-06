package pm3lib

import (
	"fmt"

	"go.bug.st/serial"
)

var (
	serialMode = serial.Mode{
		BaudRate: 115200,
	}
)

// Calls `panic` when an error is present.
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Returns a []byte as a hex string; []byte{0xff, 0xaa} = "ffaa".
func asHex(in []byte) (out string) {
	for _, byt := range in {
		out += fmt.Sprintf("%02x ", byt)
	}

	return out
}

func Connect(path string) (*Client, error) {
	client := new(Client)

	// Try to connect to the given port.
	port, err := serial.Open(path, &serialMode)
	if err != nil {
		return nil, err
	}

	// Success, give the caller a new client.
	client.port = port
	client.path = path
	return client, nil
}
