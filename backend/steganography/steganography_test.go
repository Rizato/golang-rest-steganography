package steganography

import (
	"image"
	"image/color"
	"testing"
)

func TestEmbedAndExtract(t *testing.T) {
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{R: 100, G: 150, B: 200, A: 255})
		}
	}

	testMessage := "Hello, World!"
	
	// Test embedding
	embedded, err := EmbedLsb(img, testMessage)
	if err != nil {
		t.Fatalf("Failed to embed message: %v", err)
	}

	// Test extraction
	extracted, err := ExtractLsb(embedded)
	if err != nil {
		t.Fatalf("Failed to extract message: %v", err)
	}

	if extracted != testMessage {
		t.Errorf("Extracted message doesn't match. Expected: %q, Got: %q", testMessage, extracted)
		
		// Debug: print first few bytes
		t.Logf("Original message bytes: %v", []byte(testMessage))
		t.Logf("Extracted message bytes: %v", []byte(extracted))
	}
}

func TestEmbedExtractDetailed(t *testing.T) {
	// Create a simple test image
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{R: 128, G: 128, B: 128, A: 255})
		}
	}

	testMessage := "A"
	
	// Embed message
	embedded, err := EmbedLsb(img, testMessage)
	if err != nil {
		t.Fatalf("Failed to embed: %v", err)
	}

	// Check what was actually embedded
	rgba := embedded.(*image.RGBA)
	t.Logf("Message: %q, Bytes: %v", testMessage, []byte(testMessage))
	t.Logf("Length header (4 bytes): %08b %08b %08b %08b", 0, 0, 0, 1)
	t.Logf("Message byte 'A': %08b (%d)", 'A', 'A')
	
	// Print first few pixels
	for i := 0; i < 5; i++ {
		x := i % 10
		y := i / 10
		r, g, b, _ := rgba.At(x, y).RGBA()
		r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
		t.Logf("Pixel %d: R=%08b G=%08b B=%08b (LSBs: %d%d%d)", 
			i, r8, g8, b8, r8&1, g8&1, b8&1)
	}

	// Extract and compare
	extracted, err := ExtractLsb(embedded)
	if err != nil {
		t.Fatalf("Failed to extract: %v", err)
	}

	if extracted != testMessage {
		t.Errorf("Mismatch! Expected: %q, Got: %q", testMessage, extracted)
		if len(extracted) > 0 {
			t.Logf("Extracted bytes: %v", []byte(extracted))
		}
	} else {
		t.Logf("Success! Extracted: %q", extracted)
	}
}