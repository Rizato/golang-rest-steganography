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
		image   string
		message string
		outfile string
	)

	flag.StringVar(&image, "i", "", "image to embed message in")
	flag.StringVar(&message, "m", "", "message to embed")
	flag.StringVar(&outfile, "o", "", "output file")
	flag.Parse()

	if image == "" {
		fmt.Fprintln(os.Stderr, "image is required")
		os.Exit(1)
	}

	if outfile == "" {
		fmt.Fprintln(os.Stderr, "output file is required")
		os.Exit(1)
	}

	if message == "" {
		fmt.Fprintln(os.Stderr, "output file is required")
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

	embeded, err := steganography.EmbedLsb(toEmbed, message)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	out, err := os.Create(outfile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer out.Close()
	// Use no compression and no filtering to preserve LSBs
	encoder := &png.Encoder{
		CompressionLevel: png.NoCompression,
		BufferPool:       nil,
	}
	err = encoder.Encode(out, embeded)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
