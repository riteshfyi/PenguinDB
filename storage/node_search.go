package storage

import "bytes"

func nodeLookupLE(node BNode, key []byte) uint16 {
	nkeys := node.nkeys()
	found := uint16(0)

	for i := uint16(1); i < nkeys; i++ {
		cmp := bytes.Compare(node.getKey(i), key)

		if cmp >= 0 {
			found = i
		}

		if cmp < 0 {
			break
		}
	}
	return found
}
