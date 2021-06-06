package main

import (
	"fmt"
	"log"

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
		Command: []byte{0x85, 0x03},
		NG:      false,
		Data:    []byte{0x03, 0x00, 0x00},
	}

	res, err := client.SendNGCommand(&command, true)
	log.Println(res.String())
}
