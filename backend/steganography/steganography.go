package steganography

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

var (
	ErrMessageTooLarge = errors.New("message too large for image")
	ErrInvalidImage    = errors.New("invalid image format")
	ErrNoMessage       = errors.New("no message found in image")
)

// EmbedLsb embeds a text message into a PNG image using LSB steganography
func EmbedLsb(img image.Image, message string) (image.Image, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Always convert to RGBA for consistent pixel access
	result := image.NewRGBA(bounds)
	
	// Fill with white background first to handle transparency
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			result.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}
	
	// Now draw the image over the white background
	draw.Draw(result, bounds, img, bounds.Min, draw.Over)

	// Prepare message data with 4-byte big-endian length prefix
	messageBytes := []byte(message)
	if len(messageBytes) == 0 {
		return nil, errors.New("message cannot be empty")
	}

	// Create data to embed: [4 bytes length][message bytes]
	totalData := make([]byte, 4+len(messageBytes))
	binary.BigEndian.PutUint32(totalData[0:4], uint32(len(messageBytes)))
	copy(totalData[4:], messageBytes)


	// Check capacity: we need 8 bits per data byte, 3 channels per pixel
	totalBitsNeeded := len(totalData) * 8
	availableBits := width * height * 3 // R, G, B channels

	if totalBitsNeeded > availableBits {
		return nil, ErrMessageTooLarge
	}

	// Embed bits sequentially: R0, G0, B0, R1, G1, B1, ...
	bitIndex := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if bitIndex >= totalBitsNeeded {
				break // All data embedded
			}

			// Get current pixel
			pixel := result.RGBAAt(x, y)

			// Embed in R channel
			if bitIndex < totalBitsNeeded {
				byteIndex := bitIndex / 8
				bitPosition := bitIndex % 8
				dataBit := (totalData[byteIndex] >> (7 - bitPosition)) & 1
				pixel.R = (pixel.R & 0xFE) | dataBit
				bitIndex++
			}

			// Embed in G channel
			if bitIndex < totalBitsNeeded {
				byteIndex := bitIndex / 8
				bitPosition := bitIndex % 8
				dataBit := (totalData[byteIndex] >> (7 - bitPosition)) & 1
				pixel.G = (pixel.G & 0xFE) | dataBit
				bitIndex++
			}

			// Embed in B channel
			if bitIndex < totalBitsNeeded {
				byteIndex := bitIndex / 8
				bitPosition := bitIndex % 8
				dataBit := (totalData[byteIndex] >> (7 - bitPosition)) & 1
				pixel.B = (pixel.B & 0xFE) | dataBit
				bitIndex++
			}

			result.SetRGBA(x, y, pixel)
		}
		if bitIndex >= totalBitsNeeded {
			break
		}
	}


	return result, nil
}

// ExtractLsb extracts a text message from a PNG image using LSB steganography
func ExtractLsb(img image.Image) (string, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Always convert to RGBA for consistent pixel access
	rgba := image.NewRGBA(bounds)
	
	// Fill with white background first to handle transparency
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rgba.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}
	
	// Now draw the image over the white background
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Over)
	

	// First, extract 4 bytes to get message length
	lengthBits := make([]uint8, 32) // 4 bytes * 8 bits = 32 bits
	bitIndex := 0

	for y := 0; y < height && bitIndex < 32; y++ {
		for x := 0; x < width && bitIndex < 32; x++ {
			pixel := rgba.RGBAAt(x, y)

			// Extract from R channel
			if bitIndex < 32 {
				lengthBits[bitIndex] = pixel.R & 1
				bitIndex++
			}

			// Extract from G channel  
			if bitIndex < 32 {
				lengthBits[bitIndex] = pixel.G & 1
				bitIndex++
			}

			// Extract from B channel
			if bitIndex < 32 {
				lengthBits[bitIndex] = pixel.B & 1
				bitIndex++
			}
		}
	}

	// Convert length bits to bytes
	lengthBytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		for j := 0; j < 8; j++ {
			bit := lengthBits[i*8+j]
			lengthBytes[i] = (lengthBytes[i] << 1) | bit
		}
	}

	messageLen := binary.BigEndian.Uint32(lengthBytes)

	// Validate message length
	if messageLen == 0 {
		return "", ErrNoMessage
	}

	maxPossibleLen := uint32((width*height*3 - 32) / 8) // Available bits minus length header
	if messageLen > maxPossibleLen {
		return "", fmt.Errorf("message length %d exceeds maximum possible %d", messageLen, maxPossibleLen)
	}

	// Extract message data
	totalBitsNeeded := 32 + int(messageLen)*8 // Length header + message
	messageBits := make([]uint8, int(messageLen)*8)
	bitIndex = 0
	messageBitIndex := 0

	for y := 0; y < height && bitIndex < totalBitsNeeded; y++ {
		for x := 0; x < width && bitIndex < totalBitsNeeded; x++ {
			pixel := rgba.RGBAAt(x, y)

			// Extract from R channel
			if bitIndex >= 32 && messageBitIndex < len(messageBits) {
				messageBits[messageBitIndex] = pixel.R & 1
				messageBitIndex++
			}
			bitIndex++

			// Extract from G channel
			if bitIndex >= 32 && messageBitIndex < len(messageBits) {
				messageBits[messageBitIndex] = pixel.G & 1
				messageBitIndex++
			}
			bitIndex++

			// Extract from B channel
			if bitIndex >= 32 && messageBitIndex < len(messageBits) {
				messageBits[messageBitIndex] = pixel.B & 1
				messageBitIndex++
			}
			bitIndex++
		}
	}

	// Convert message bits to bytes
	messageBytes := make([]byte, messageLen)
	for i := 0; i < int(messageLen); i++ {
		for j := 0; j < 8; j++ {
			if i*8+j < len(messageBits) {
				bit := messageBits[i*8+j]
				messageBytes[i] = (messageBytes[i] << 1) | bit
			}
		}
	}

	return string(messageBytes), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// EncodeImageToPNG encodes an image to PNG format bytes
func EncodeImageToPNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}
	return buf.Bytes(), nil
}

// DecodeImageFromPNG decodes PNG bytes to an image
func DecodeImageFromPNG(data []byte) (image.Image, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG: %w", err)
	}
	return img, nil
}
