package pm3lib

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
)

var (
	clientMagicPreamble = []byte{0x50, 0x4D, 0x33, 0x61}
	deviceMagicPreamble = []byte{0x50, 0x4D, 0x33, 0x62}

	clientMagicChecksum = []byte{0x61, 0x33}
	deviceMagicChecksum = []byte{0x62, 0x33}

	minDeviceResponseSize = 8
)

type NGCommand struct {
	Command []byte // 2 bytes
	Data    []byte // 0-uint15 bytes.
	NG      bool   // If in NG mode.
}

func (n *NGCommand) serialize() ([]byte, error) {
	// Begins with the preamble.
	out := clientMagicPreamble

	// Convert the length of `n.Data` into two little-endian bytes.
	lengthBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(lengthBytes, uint16(len(n.Data)))

	// If we are in NG mode, the last bit must be whether or not NG is enabled.
	if n.NG {
		lengthBytes[1] |= (1 << 7)
	}

	// Now we can append the properly-formatted length bytes.
	out = append(out, lengthBytes...)

	// The two-byte command next...
	if len(n.Command) != 2 {
		return nil, errors.New("n.Command is not 2 bytes long")
	}
	out = append(out, n.Command...)

	// Output data next.
	out = append(out, n.Data...)

	// Finally, the fake CRC.
	out = append(out, clientMagicChecksum...)
	return out, nil
}

type NGResponse struct {
	Command []byte // 2 bytes
	Data    []byte // 0-uint15 bytes.
	DataLen uint16

	NG     bool // If in NG mode.
	Status int16

	loaded bool
}

func (r *NGResponse) load(res []byte) error {
	log.Println("Processing a response:", asHex(res))

	// Sanity check on how long this is.
	if len(res) < minDeviceResponseSize {
		return errors.New("Response is too small")
	}

	// Ensure the preamble is part of the response.
	if bytes.Compare(res[0:4], deviceMagicPreamble) != 0 {
		return errors.New("Preamble is missing from response")
	}

	// Ensure the fake checksum is at the end of the response.
	if bytes.Compare(res[len(res)-2:], deviceMagicChecksum) != 0 {
		return errors.New("Fake checksum is missing from response")
	}

	// Check if the response is NG and clear the bit.
	r.NG = res[5]&(1<<7) == 1
	res[5] &= ^(byte(1) << 7)

	// Set the data length.
	r.DataLen = binary.LittleEndian.Uint16(res[4:6])

	// Set the status.
	r.Status = int16(binary.LittleEndian.Uint16(res[6:8]))

	// Set the command.
	r.Command = res[8:10]

	// Copy the data, if any.
	if r.DataLen > 0 {
		r.Data = res[10 : len(res)-2]
	}

	r.loaded = true
	return nil
}

func (r *NGResponse) String() (output string) {
	if !r.loaded {
		panic("Tried to stringify a NGResponse before it was loaded")
	}

	output += "\n"
	output += fmt.Sprintf("Status   : %d\n", r.Status)
	output += fmt.Sprintf("Command  : %s\n", asHex(r.Command))
	output += fmt.Sprintf("NG       : %t\n", r.NG)
	output += fmt.Sprintf("Data Len : %d\n", r.DataLen)
	output += fmt.Sprintf("Data     : %s\n", asHex(r.Data))

	return output
}
