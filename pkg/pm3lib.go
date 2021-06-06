package pm3lib

import (
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
