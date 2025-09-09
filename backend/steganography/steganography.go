package steganography

import (
	"errors"
	"image"
	"image/png"
)

var (
	Failed = errors.New("failed")
)

func EmbedLSB(image image.Image, messsage string) (image.Image, error) {
	
}

func ExtractLsb(image image.Image) (string, error) {

}
