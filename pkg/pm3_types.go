package pm3lib

import (
	"encoding/binary"
	"errors"
)

var (
	clientMagicPreamble = []byte{0x50, 0x4D, 0x33, 0x61}
	deviceMagicPreamble = []byte{0x50, 0x4D, 0x33, 0x62}

	clientMagicChecksum = []byte{0x61, 0x33}
	deviceMagicChecksum = []byte{0x62, 0x33}
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
