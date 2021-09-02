package pdutext

import (
	smpppdutext "github.com/fiorix/go-smpp/smpp/pdu/pdutext"
)

// Raw text codec, no encoding.
type Raw []byte

// Type implements the Codec interface.
func (s Raw) Type() smpppdutext.DataCoding {

	//return smpppdutext.ISO88595Type
	var t smpppdutext.DataCoding = 0x48
	return t
}

// Encode raw text.
func (s Raw) Encode() []byte {
	return s
}

// Decode raw text.
func (s Raw) Decode() []byte {
	return s
}
