package main

import (
	"fmt"

	pm3lib "github.com/iangcarroll/pm3lib/pkg"
)

const (
	devicePath = "/dev/tty.usbmodemiceman1"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Returns a []byte as a hex string; []byte{0xff, 0xaa} = "ffaa".
func asHex(in []byte) (out string) {
	for _, byt := range in {
		out += fmt.Sprintf("%02x", byt)
	}

	return out
}

func main() {
	client, err := pm3lib.Connect(devicePath)
	check(err)

	command := pm3lib.NGCommand{
		Command: []byte{0x30, 0x04},
		NG:      true,
		Data:    []byte{},
	}
	client.SendNGCommand(&command, false)
}
