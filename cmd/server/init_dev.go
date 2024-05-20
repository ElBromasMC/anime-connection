//go:build dev

package main

import (
	"alc/model/cart"
	"encoding/gob"
	"net/http"
)

func init() {
	// Live reload
	http.Get("http://localhost:8020")

	gob.Register([]byte{})
	gob.Register([]cart.ItemRequest{})
}
