package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Fprintf(os.Stdout, "%s is a dir: %t\n", file.Name(), file.IsDir())
	}
}
