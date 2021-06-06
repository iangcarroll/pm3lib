package pm3lib

import (
	"fmt"
	"sync"

	"go.bug.st/serial"
)

type Client struct {
	path  string
	port  serial.Port
	mutex sync.Mutex
}

const (
	responseBufferLength = 1024
)

func (c *Client) SendNGCommand(n *NGCommand, hasResponse bool) ([]byte, error) {
	serialized, err := n.serialize()
	if err != nil {
		return nil, err
	}

	return c.transmit(serialized, hasResponse)
}

// Sends `payload` to the device and checks for a response.
func (c *Client) transmit(payload []byte, hasResponse bool) ([]byte, error) {
	// Write the payload to the serial console.
	n, err := c.port.Write(payload)

	// Check if the write failed.
	if err != nil {
		return []byte{}, err
	}

	// Check if we underwrote.
	if len(payload) != n {
		return []byte{}, fmt.Errorf("tried to write %d bytes but only wrote %d", len(payload), n)
	}

	// If the command has no response, we're done.
	if !hasResponse {
		return []byte{}, nil
	}

	// Allocate a buffer for the response and try to read.
	resBuffer := make([]byte, responseBufferLength)

	// Try to read out the bytes (blocking).
	readBytes, err := c.port.Read(resBuffer)
	if err != nil {
		return []byte{}, err
	}

	// If we have response bytes, return them.
	if readBytes > 0 {
		return resBuffer[:readBytes], nil
	}

	// Return an empty response.
	return []byte{}, nil
}
