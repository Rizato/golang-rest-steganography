package filestype

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

var DefaultMaxFilesize int64 = 64 * 1024 * 1024            // 64 meg
var DefaultMimeTypes = []string{"image/jpeg", "image/png"} // just jpeg and png

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
		return FileTooLargeError{}
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
		return InvalidMimeType{}
	}

	return nil
}
