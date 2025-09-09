package filestype

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

// Config
var (
	DefaultMaxFilesize int64 = 64 * 1024 * 1024      // 64 meg
	DefaultMimeTypes         = []string{"image/png"} // just png for now
	FileTooLargeError        = errors.New("file too large")
	InvalidMimetype          = errors.New("invalid mimetype")
)

// DefaultValidator Global default validator
var DefaultValidator = NewFileValidator(DefaultMaxFilesize, DefaultMimeTypes...)

// Validate Helper validator that uses the global default validator
func Validate(fileheader *multipart.FileHeader) error {
	return DefaultValidator.Validate(fileheader)
}

func GetMimetype(reader io.Reader) (string, error) {
	data := make([]byte, 2048)
	s, err := reader.Read(data)
	if err != nil {
		return "", err
	}
	return http.DetectContentType(data[:s]), nil
}

func ValidateMimetype(mimetype string) error {
	return DefaultValidator.ValidateMimetype(mimetype)
}

// FileValidator Validate files given the allowed configuration
type FileValidator struct {
	MaxFilesize  int64
	AllowedTypes map[string]bool
}

func NewFileValidator(maxFilesize int64, mimetypes ...string) FileValidator {
	mt := make(map[string]bool)
	for _, m := range mimetypes {
		mt[m] = true
	}
	return FileValidator{
		MaxFilesize:  maxFilesize,
		AllowedTypes: mt,
	}
}

func (v *FileValidator) Validate(fileheader *multipart.FileHeader) error {
	if fileheader.Size > v.MaxFilesize {
		return FileTooLargeError
	}

	mimetypes := fileheader.Header["Content-Type"]

	for _, mt := range mimetypes {
		err := v.ValidateMimetype(mt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *FileValidator) ValidateMimetype(mimetype string) error {
	if !v.AllowedTypes[mimetype] {
		return InvalidMimetype
	}
	return nil
}
