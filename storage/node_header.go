package storage

import "encoding/binary"

const (
	HEADER = uint16(4)
)

func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node.data)
}

func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node.data[2:])
}

func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node.data[0:], btype)
	binary.LittleEndian.PutUint16(node.data[2:], nkeys)
}
