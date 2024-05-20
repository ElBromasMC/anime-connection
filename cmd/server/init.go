package main

import (
	"alc/model/cart"
	"encoding/gob"
)

func init() {
	gob.Register([]byte{})
	gob.Register([]cart.ItemRequest{})
}
