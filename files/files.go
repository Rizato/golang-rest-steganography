package filestype

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

var (
	DefaultMaxFilesize int64 = 64 * 1024 * 1024                    // 64 meg
	DefaultMimeTypes         = []string{"image/jpeg", "image/png"} // just jpeg and png
	FileTooLargeError        = errors.New("file too large")
	InvalidMimeType          = errors.New("invalid mime type")
)

type MultiPartFileSaver struct {
	r              *http.Request
	uploadFilename string
	path           string
	name           string
	validator      FileValidator
}

func NewMultiPartFileSaver(r *http.Request, uploadFilename string, path string, name string, validator FileValidator) *MultiPartFileSaver {
	return &MultiPartFileSaver{
		r,
		uploadFilename,
		path,
		name,
		validator,
	}
}

func (s *MultiPartFileSaver) SaveFile() (string, error) {
	err := s.r.ParseMultipartForm(32 << 20)
	if err != nil {
		return "", err
	}

	// Get the uploaded from the form
	uploaded, uploadedHeader, err := s.r.FormFile(s.uploadFilename)
	if err != nil {
		return "", err
	}
	defer uploaded.Close()

	// Use LimitReader to stop reading beyond file size, + 1 so I can detect oversized files with spoofed values
	limited := io.LimitReader(uploaded, DefaultMaxFilesize+1)

	// Validate the size and mime types (though user supplied and can be manipulated)
	err = s.validator.ValidateImage(uploadedHeader)
	if err != nil {
		return "", err
	}

	// Write to an specified dir
	file, err := os.Create(s.path + s.name)
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileSize, err := io.Copy(file, limited)
	if err != nil {
		return "", err
	}
	// Don't trust client set size
	if fileSize > s.validator.maxFilesize {
		return "", FileTooLargeError
	}
	return file.Name(), nil
}

type FileValidator struct {
	maxFilesize int64
	mimetypes   map[string]bool
}

func NewFileValidator(maxFilesize int64, mimetypes ...string) FileValidator {
	mt := make(map[string]bool)
	for _, m := range mimetypes {
		mt[m] = true
	}
	return FileValidator{
		maxFilesize: maxFilesize,
		mimetypes:   mt,
	}
}

func (v *FileValidator) ValidateImage(fileheader *multipart.FileHeader) error {
	if fileheader.Size > v.maxFilesize {
		return FileTooLargeError
	}

	mimetypes := fileheader.Header["Content-Type"]
	mimeMatch := false
	for _, mt := range mimetypes {
		if v.mimetypes[mt] {
			mimeMatch = true
			break
		}
	}

	if !mimeMatch {
		return InvalidMimeType
	}

	return nil
}
