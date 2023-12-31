package main

import (
	"errors"
	"log"
)

func logtest() {
	log.Printf("测试")
	log.Println("测试")
	err := errors.New("test")
	fatalIfErr(err, "delete cert")
}

func fatalIfErr(err error, msg string) {
	if err != nil {
		log.Fatalf("ERROR: %s: %s", msg, err)
	}
}
