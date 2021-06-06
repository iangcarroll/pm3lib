package pm3lib

import (
	"errors"
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

func (c *Client) SendNGCommand(n *NGCommand, hasResponse bool) (*NGResponse, error) {
	serialized, err := n.serialize()
	if err != nil {
		return nil, err
	}

	return c.transmit(serialized, hasResponse)
}

// Sends `payload` to the device and checks for a response.
func (c *Client) transmit(payload []byte, hasResponse bool) (*NGResponse, error) {
	// Write the payload to the serial console.
	n, err := c.port.Write(payload)

	// Check if the write failed.
	if err != nil {
		return nil, err
	}

	// Check if we underwrote.
	if len(payload) != n {
		return nil, fmt.Errorf("tried to write %d bytes but only wrote %d", len(payload), n)
	}

	// If the command has no response, we're done.
	if !hasResponse {
		return nil, nil
	}

	// Allocate a buffer for the response and try to read.
	resBuffer := make([]byte, responseBufferLength)

	// Try to read out the bytes (blocking).
	readBytes, err := c.port.Read(resBuffer)
	if err != nil {
		return nil, err
	}

	// If we have a response, return it.
	if readBytes > 0 {
		response := new(NGResponse)
		return response, response.load(resBuffer[:readBytes])
	}

	// Return an empty response.
	return nil, errors.New("no data from read")
}
