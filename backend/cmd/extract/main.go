package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"
	"steg/steganography"
)

func main() {
	var (
		image string
	)

	flag.StringVar(&image, "i", "", "image to extract a message from")
	flag.Parse()

	if image == "" {
		fmt.Fprintln(os.Stderr, "image is required")
		os.Exit(1)
	}

	file, err := os.Open(image)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()

	// Decode file to a png
	toEmbed, err := png.Decode(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	message, err := steganography.ExtractLsb(toEmbed)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(message)
}
