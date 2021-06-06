## pm3lib
`pm3lib` is a Go package for connecting to and interfacing with Proxmark3 devices. It allows you to build custom host applications that call Proxmark commands implemented by its firmware. Using these bindings, you could implement your own Proxmark client, scripts for a connected device, etc.

### Getting started
In your own Go application, using `pm3lib` is pretty easy. For example, to connect to a Proxmark3 and drop the reader field:

```go
import pm3lib "github.com/iangcarroll/pm3lib/pkg"

func main() {
	client, err := pm3lib.Connect("/dev/tty.usbmodemiceman1")
	check(err)

	client.DropField()
}
```

Most Proxmark3 commands have not been implemented directly in the client, so you will need to construct them yourself. In the RRG firmware, commands are documented in [pm3_cmd.h](https://github.com/RfidResearchGroup/proxmark3/blob/master/include/pm3_cmd.h#L654). The DropField command can be called manually like this:
```go
import pm3lib "github.com/iangcarroll/pm3lib/pkg"

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
```

Similarly, for commands which have response data, you can access a structured `NGResponse`. This example calls `CMD_HF_ISO14443A_READER`. The Proxmark also emits unsolicited responses for a variety of reasons, i.e. debug logging. You will have to handle (or discard) these unrelated responses yourself.

```go
import pm3lib "github.com/iangcarroll/pm3lib/pkg"

func main() {
	client, err := pm3lib.Connect(devicePath)
	check(err)

	command := pm3lib.NGCommand{
		Command: []byte{0x85, 0x03},
		NG:      true,
		Data:    []byte{0x03, 0x00, 0x00},
	}

	err = client.SendNGCommand(&command)
	check(err)

	for {
		res, err := client.ReceiveNGResponse()
		check(err)
		
		log.Println(res.String())
	}
}
```