package main

import (
	"log"

	"github.com/online-bnsp/backend/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
