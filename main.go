package main

import (
	"log"
)

func main() {
	//test
	store, err := NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewApiServer(":8080", store)
	server.Run()
}
