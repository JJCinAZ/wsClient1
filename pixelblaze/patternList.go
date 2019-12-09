package pixelblaze

import (
	"bytes"
	"fmt"
)

type Program struct {
	ID   string
	Name string
}

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (p *transProp) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &p.a, &p.b)
	return err
}
