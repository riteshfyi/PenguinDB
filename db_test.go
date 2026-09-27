package main

import (
	"PenguinDB/storage"
	"fmt"
	"testing"
)

func TestDB(t *testing.T) {
	db := storage.KV{Path: "./database/test_file.kv"}
	err := db.Open()
	if err != nil {
		fmt.Print(err)
	}
	keyString1 := "kaguyaSama"
	valString1 := "LoveisWar"

	key1 := []byte(keyString1)
	val1 := []byte(valString1)
	db.Set(key1, val1)

	ks2 := "Tora"
	vs2 := "Dora"

	key2 := []byte(ks2)
	val2 := []byte(vs2)

	db.Set(key2, val2)

	setVal1, ok1 := db.Get(key1)

	if !ok1 {
		t.Fatal("Set Key1 not found")
	} else {
		fmt.Println(string(setVal1))
	}

	setVal2, ok2 := db.Get(key2)

	if !ok2 {
		t.Fatal("Set Key2 not found")
	} else {
		fmt.Println(string(setVal2))
	}

	del1, err := db.Del(key1)

	if err != nil {
		t.Fatal("Unable to delete key1")
	} else if !del1 {
		t.Fatal("Set Key not found for Delete")
	}

	db.Close()
}
