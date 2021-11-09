package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var inputPath string
	flag.StringVar(&inputPath, "InputFile", "none", "Specify the binary hibp file")
	flag.Parse()

	if inputPath == "none" {
		fmt.Println(os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	inputFile, err := os.Open(inputPath)
	handleErr(err)
	defer inputFile.Close()

	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	handleErr(err)
	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	handleErr(err)
	fmt.Println()
	if err = term.Restore(int(os.Stdin.Fd()), state); err != nil {
		log.Fatal(err)
	}

	digest := sha1.New()
	size, err := digest.Write(bytePassword)
	handleErr(err)
	if size != len(bytePassword) {
		log.Fatal("Could not hash the password. Aborting.")
	}
	passwordHash := digest.Sum(nil)

	fmt.Println("Your Password Hash is:", hex.EncodeToString(passwordHash))

	hashSize := int64(len(passwordHash))
	buffer := make([]byte, hashSize)

	fileInfo, err := inputFile.Stat()
	handleErr(err)

	lower := int64(0)
	upper := (fileInfo.Size() / hashSize) - 1

	for lower <= upper {
		position := lower + ((upper - lower) / 2)

		inputFile.Seek(position*hashSize, 0)
		readCount, err := inputFile.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		if readCount != int(hashSize) {
			log.Fatal("Did not get a valid hash. Aborting.")
		}

		comparison := bytes.Compare(passwordHash, buffer)

		if comparison == 0 {
			fmt.Println("Password pwned!")
			os.Exit(0)
		}

		if comparison < 0 {
			upper = position - 1
		} else if comparison > 0 {
			lower = position + 1
		}
	}

	fmt.Println("Password safe :)")
	os.Exit(0)
}
