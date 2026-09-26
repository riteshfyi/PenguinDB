package main

import (
	"PenguinDB/storage"
	"fmt"
	"testing"
)

func TestBTree(t *testing.T) {

	c := storage.NewC()
	key := "kaguya"
	val := "sama"
	c.Add(key, val)
	storedVal, ok := c.Get(key)

	if !ok {
		panic("key not present in the tree, however it was stored")
	}

	if storedVal != val {
		t.Fatal("stored key doesn't match fetched key")
	} else {
		fmt.Println(storedVal)
	}

	// c.Del(key)

	// _, ok1 := c.Get(key)

	// if ok1 {
	// 	panic("key still present after deletion")

	// }

	key2 := "miyuki"
	val2 := "shirogane"

	c.Add(key2, val2)

	storedVal2, ok := c.Get(key2)

	if storedVal2 != val2 {
		t.Fatal("stored key doesn't match fetched key")
	} else {
		fmt.Println(storedVal2)
	}

	c.Del(key)

	_, ok1 := c.Get(key)

	if ok1 {
		panic("key still present after deletion")

	}

	c.Del(key2)

	_, ok2 := c.Get(key2)

	if ok2 {
		t.Fatal("key2 still present after deletion")

	}
}
