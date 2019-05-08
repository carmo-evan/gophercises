package util

import (
	"encoding/binary"
	"log"
)

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
