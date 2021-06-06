package main

import (
	"bytes"
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
		Command: []byte{0xAD, 0xDE},
		NG:      true,
		Data:    []byte{},
	}

	err = client.SendNGCommand(&command)
	log.Println(err)

	for {
		res, err := client.ReceiveNGResponse()
		check(err)

		log.Println(res.String())

		if bytes.Compare(res.Command, []byte{0x00, 0x01}) == 0 {
			log.Println("DEBUG LOG!", string(res.Command[2:]))
		}
	}
}
