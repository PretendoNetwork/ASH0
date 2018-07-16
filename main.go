package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func decompressAsh(ashFile string) []byte {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(dir+"\\AshUtils.exe", "decompress", "test.ash")
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("%s", out)
		log.Fatal(err)
	}

	return out
}

func main() {
	out := decompressAsh("test.ash")
	fmt.Println("Outputting to file...")
	ioutil.WriteFile("out.ash", out, 0644)
}
