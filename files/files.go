package filestype

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type FileTooLargeError struct{}

func (e FileTooLargeError) Error() string {
	return "File too large"
}

type InvalidMimeType struct{}

func (e InvalidMimeType) Error() string {
	return "Invalid Mime type"
}

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
	_, err = io.Copy(file, uploaded)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

type FileValidator struct {
	maxFileSize    int64
	validMimeTypes map[string]bool
}

func (v *FileValidator) ValidateImage(fileheader *multipart.FileHeader) error {
	if fileheader.Size > v.maxFileSize {
		return FileTooLargeError{}
	}

	mimetypes := fileheader.Header["Content-Type"]
	validMimeType := false
	for _, mt := range mimetypes {
		if v.validMimeTypes[mt] {
			validMimeType = true
			break
		}
	}

	if !validMimeType {
		return InvalidMimeType{}
	}

	return nil
}
