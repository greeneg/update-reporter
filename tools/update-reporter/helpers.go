package main

import (
	"log"
)

func fatalCheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
