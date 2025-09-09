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
	// Convert image to RGBA for manipulation
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// Prepare message with length header
	messageBytes := []byte(message)
	messageLen := uint32(len(messageBytes))

	// Check if image has enough capacity (3 bits per pixel for RGB)
	maxCapacity := (bounds.Dx() * bounds.Dy() * 3) / 8 // bits / 8 = bytes
	totalSize := 4 + len(messageBytes)                 // 4 bytes for length + message

	if totalSize > maxCapacity {
		return nil, ErrMessageTooLarge
	}

	// Create data buffer with length prefix
	data := make([]byte, 0, totalSize)
	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, messageLen)
	data = append(data, lengthBuf...)
	data = append(data, messageBytes...)

	// Embed data into image
	dataIndex := 0
	bitIndex := 0

	for y := bounds.Min.Y; y < bounds.Max.Y && dataIndex < len(data); y++ {
		for x := bounds.Min.X; x < bounds.Max.X && dataIndex < len(data); x++ {
			r, g, b, a := rgba.At(x, y).RGBA()

			// Convert to 8-bit values
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			// Embed bits in R channel
			if dataIndex < len(data) {
				bit := (data[dataIndex] >> (7 - bitIndex)) & 1
				r8 = (r8 & 0xFE) | bit
				bitIndex++
				if bitIndex == 8 {
					bitIndex = 0
					dataIndex++
				}
			}

			// Embed bits in G channel
			if dataIndex < len(data) {
				bit := (data[dataIndex] >> (7 - bitIndex)) & 1
				g8 = (g8 & 0xFE) | bit
				bitIndex++
				if bitIndex == 8 {
					bitIndex = 0
					dataIndex++
				}
			}

			// Embed bits in B channel
			if dataIndex < len(data) {
				bit := (data[dataIndex] >> (7 - bitIndex)) & 1
				b8 = (b8 & 0xFE) | bit
				bitIndex++
				if bitIndex == 8 {
					bitIndex = 0
					dataIndex++
				}
			}

			rgba.SetRGBA(x, y, color.RGBA{R: r8, G: g8, B: b8, A: a8})
		}
	}

	return rgba, nil
}

// ExtractLSB extracts a text message from a PNG image using LSB steganography
func ExtractLSB(img image.Image) (string, error) {
	bounds := img.Bounds()

	// Extract all data bits first
	maxBytes := (bounds.Dx() * bounds.Dy() * 3) / 8
	data := make([]byte, 0, maxBytes)

	bitIndex := 0
	currentByte := byte(0)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// Extract from R channel
			r8 := uint8(r >> 8)
			bit := r8 & 1
			currentByte = (currentByte << 1) | bit
			bitIndex++
			if bitIndex == 8 {
				data = append(data, currentByte)
				bitIndex = 0
				currentByte = 0
			}

			// Extract from G channel
			g8 := uint8(g >> 8)
			bit = g8 & 1
			currentByte = (currentByte << 1) | bit
			bitIndex++
			if bitIndex == 8 {
				data = append(data, currentByte)
				bitIndex = 0
				currentByte = 0
			}

			// Extract from B channel
			b8 := uint8(b >> 8)
			bit = b8 & 1
			currentByte = (currentByte << 1) | bit
			bitIndex++
			if bitIndex == 8 {
				data = append(data, currentByte)
				bitIndex = 0
				currentByte = 0
			}

			// Stop if we have enough data for length + max possible message
			if len(data) >= 4 {
				messageLen := binary.BigEndian.Uint32(data[:4])
				if messageLen > 0 && messageLen < uint32(maxBytes) && len(data) >= int(4+messageLen) {
					return string(data[4 : 4+messageLen]), nil
				}
			}
		}
	}

	// Check if we have at least the length header
	if len(data) < 4 {
		return "", ErrNoMessage
	}

	messageLen := binary.BigEndian.Uint32(data[:4])
	if messageLen == 0 || messageLen > uint32(len(data)-4) {
		return "", ErrNoMessage
	}

	return string(data[4 : 4+messageLen]), nil
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
