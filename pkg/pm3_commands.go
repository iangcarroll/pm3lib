package pm3lib

var (
	dropFieldCommand = NGCommand{
		Command: []byte{0x30, 0x04},
		NG:      true,
		Data:    []byte{},
	}
)

func (c *Client) DropField() error {
	return c.SendNGCommand(&dropFieldCommand)
}
