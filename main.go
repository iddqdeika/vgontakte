package main

import (
	"sync"
	"time"
	"vgontakte/storage/localstorage"
)

func main() {

	m := sync.RWMutex{}

	storage := localstorage.GetLocalStorage("testdb")

	err := storage.Update("path.key", "value1")

	if err != nil {
		panic(err)
	}

	val, err := storage.Get("path.key")
	if err != nil {
		panic(err)
	}
	println(string(val))

}
